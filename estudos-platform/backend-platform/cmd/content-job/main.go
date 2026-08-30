package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/usecase"
	"github.com/thiago-tertuliano/estudos-platform/internal/infrastructure/config"
	"github.com/thiago-tertuliano/estudos-platform/internal/infrastructure/content"
	"github.com/thiago-tertuliano/estudos-platform/internal/infrastructure/persistence/postgres"
	pgrepo "github.com/thiago-tertuliano/estudos-platform/internal/infrastructure/persistence/postgres/repository"
)

func main() {
	dryRun := flag.Bool("dry-run", false, "valida e imprime; não grava")
	autorEmail := flag.String("autor-email", "autor.seed@estudos.local", "dono dos artigos")
	skipCatalogo := flag.Bool("skip-catalogo", false, "não lê Courses.md")
	skipAulaGo := flag.Bool("skip-aula-go", false, "não enriquece go-basico")
	flag.Parse()

	if *skipCatalogo && *skipAulaGo {
		log.Fatal("nada para importar: --skip-catalogo e --skip-aula-go juntos")
	}

	root, err := encontrarRaiz()
	if err != nil {
		log.Fatal(err)
	}
	_ = godotenv.Load()
	_ = godotenv.Load(filepath.Join(root, "backend-platform", ".env"))

	plano, err := montarPlano(root, !*skipCatalogo, !*skipAulaGo)
	if err != nil {
		log.Fatal(err)
	}
	dumpPlano(root, plano)

	if *dryRun {
		uc := usecase.NewImportarConteudo(nil, nil, nil)
		rel, err := uc.Execute(context.Background(), *autorEmail, plano, true)
		if err != nil {
			log.Fatal(err)
		}
		for _, t := range plano.Trilhas {
			fmt.Printf("trilha %s (%d módulos)\n", t.Slug, len(t.Modulos))
		}
		printRel(rel)
		return
	}

	cfg := config.LoadDB()
	pool := postgres.NewConnection(cfg)
	defer pool.Close()

	uc := usecase.NewImportarConteudo(
		pgrepo.NewUsuarioRepoPG(pool),
		pgrepo.NewTrilhaRepoPG(pool),
		pgrepo.NewArtigoRepoPG(pool),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	rel, err := uc.Execute(ctx, *autorEmail, plano, false)
	if err != nil {
		log.Fatal(err)
	}
	printRel(rel)
}

func montarPlano(root string, catalogo, aulaGo bool) (*dto.PlanoImportacao, error) {
	plano := &dto.PlanoImportacao{}
	if catalogo {
		raw, err := os.ReadFile(filepath.Join(root, "fontes", "Courses.md"))
		if err != nil {
			return nil, fmt.Errorf("Courses.md: %w", err)
		}
		cat, err := content.ParseCourses(string(raw))
		if err != nil {
			return nil, fmt.Errorf("parse catálogo: %w", err)
		}
		plano.Trilhas = append(plano.Trilhas, cat.Trilhas...)
		plano.Avisos = append(plano.Avisos, cat.Avisos...)
	}
	if aulaGo {
		raw, err := os.ReadFile(filepath.Join(root, "fontes", "AULA-CODE-REVIEW-GO-SPRINT.md"))
		if err != nil {
			return nil, fmt.Errorf("aula Go: %w", err)
		}
		goTrilha, err := content.ParseAulaGo(string(raw))
		if err != nil {
			return nil, fmt.Errorf("parse aula Go: %w", err)
		}
		plano.Trilhas = append(plano.Trilhas, *goTrilha)
	}
	return plano, nil
}

func dumpPlano(root string, plano *dto.PlanoImportacao) {
	dumpDir := filepath.Join(root, "content", "generated")
	if err := os.MkdirAll(dumpDir, 0o755); err != nil {
		return
	}
	b, err := json.MarshalIndent(plano, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(filepath.Join(dumpDir, "plano.json"), b, 0o644)
}

func printRel(rel *dto.RelatorioImportacao) {
	fmt.Printf("dry_run=%v trilhas=%d (novas=%d) artigos=%d (novos=%d)\n",
		rel.DryRun, rel.TrilhasOK, rel.TrilhasCriadas, rel.ArtigosOK, rel.ArtigosCriados)
	for _, a := range rel.Avisos {
		fmt.Println("aviso:", a)
	}
}

func encontrarRaiz() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := wd
	for i := 0; i < 8; i++ {
		if fileExists(filepath.Join(dir, "fontes", "Courses.md")) {
			return dir, nil
		}
		if fileExists(filepath.Join(dir, "estudos-platform", "fontes", "Courses.md")) {
			return filepath.Join(dir, "estudos-platform"), nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("não achei fontes/Courses.md a partir de %s", wd)
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}
