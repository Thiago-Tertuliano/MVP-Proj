import type { Meta, StoryObj } from "@storybook/react";

import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Breadcrumb, BreadcrumbItem, BreadcrumbLink, BreadcrumbList, BreadcrumbPage, BreadcrumbSeparator } from "@/components/ui/breadcrumb";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Separator } from "@/components/ui/separator";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";

/** Estrutura e navegação: Card, Breadcrumb, Tabs, Dialog, Tooltip, Avatar, Separator. */
const meta = {
  title: "UI/Estrutura",
} satisfies Meta;

export default meta;
type Story = StoryObj<typeof meta>;

export const Cartao: Story = {
  render: () => (
    <Card className="max-w-sm">
      <CardHeader>
        <CardTitle>Go Básico</CardTitle>
        <CardDescription>Fundamentos de Go para ler o backend.</CardDescription>
      </CardHeader>
      <CardContent>
        <p className="text-sm text-muted-foreground">2 módulos · 3 artigos</p>
      </CardContent>
      <CardFooter>
        <Button size="sm">Abrir trilha</Button>
      </CardFooter>
    </Card>
  ),
};

export const CartaoEscuro: Story = { ...Cartao, globals: { theme: "dark" } };

export const Migalhas: Story = {
  render: () => (
    <Breadcrumb>
      <BreadcrumbList>
        <BreadcrumbItem>
          <BreadcrumbLink href="/">Trilhas</BreadcrumbLink>
        </BreadcrumbItem>
        <BreadcrumbSeparator />
        <BreadcrumbItem>
          <BreadcrumbLink href="/trilhas/go-basico">Go Básico</BreadcrumbLink>
        </BreadcrumbItem>
        <BreadcrumbSeparator />
        <BreadcrumbItem>
          <BreadcrumbPage>Pacotes em Go</BreadcrumbPage>
        </BreadcrumbItem>
      </BreadcrumbList>
    </Breadcrumb>
  ),
};

export const Abas: Story = {
  render: () => (
    <Tabs defaultValue="artigo" className="max-w-md">
      <TabsList>
        <TabsTrigger value="artigo">Artigo</TabsTrigger>
        <TabsTrigger value="notas">Anotações</TabsTrigger>
        <TabsTrigger value="quiz" disabled>
          Quiz
        </TabsTrigger>
      </TabsList>
      <TabsContent value="artigo">Conteúdo do artigo.</TabsContent>
      <TabsContent value="notas">Suas anotações.</TabsContent>
    </Tabs>
  ),
};

export const Modal: Story = {
  render: () => (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="outline">Abrir diálogo</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Sair da conta?</DialogTitle>
          <DialogDescription>Seu progresso fica salvo. Você pode entrar de novo quando quiser.</DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline">Cancelar</Button>
          <Button>Sair</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  ),
};

export const Dica: Story = {
  render: () => (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger asChild>
          <Button variant="outline">Passe o mouse</Button>
        </TooltipTrigger>
        <TooltipContent>Marca este artigo como lido</TooltipContent>
      </Tooltip>
    </TooltipProvider>
  ),
};

export const Avatares: Story = {
  render: () => (
    <div className="flex items-center gap-3">
      <Avatar>
        <AvatarFallback>TM</AvatarFallback>
      </Avatar>
      <Avatar className="h-12 w-12">
        <AvatarFallback>AB</AvatarFallback>
      </Avatar>
    </div>
  ),
};

export const Divisores: Story = {
  render: () => (
    <div className="max-w-xs space-y-3">
      <p className="text-sm">Acima</p>
      <Separator />
      <p className="text-sm">Abaixo</p>
    </div>
  ),
};
