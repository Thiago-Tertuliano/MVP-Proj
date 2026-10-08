import Link from "next/link";
import { FileQuestion } from "lucide-react";

import { Button } from "@/components/ui/button";
import { EmptyState } from "@/components/ui/empty-state";

export default function NotFound() {
  return (
    <EmptyState
      icon={<FileQuestion />}
      headingLevel="h1"
      title="Página não encontrada"
      description="O conteúdo que você procura não existe ou ainda não foi publicado."
      action={
        <Button asChild>
          <Link href="/">Voltar para as trilhas</Link>
        </Button>
      }
    />
  );
}
