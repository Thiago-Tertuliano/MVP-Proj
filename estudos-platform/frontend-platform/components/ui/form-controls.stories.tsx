import type { Meta, StoryObj } from "@storybook/react";

import { Checkbox } from "@/components/ui/checkbox";
import { FormField } from "@/components/ui/form-field";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { Textarea } from "@/components/ui/textarea";

/** Controles de formulário: Input, Textarea, Select, Checkbox, Radio, Switch e o wrapper FormField. */
const meta = {
  title: "UI/Formulário",
  parameters: { layout: "padded" },
} satisfies Meta;

export default meta;
type Story = StoryObj<typeof meta>;

export const InputPadrao: Story = {
  name: "Input / padrão",
  render: () => (
    <FormField label="E-mail" helper="Usaremos só para o login.">
      {(c) => <Input type="email" placeholder="nome@exemplo.com" {...c} />}
    </FormField>
  ),
};

export const InputHover: Story = {
  name: "Input / hover",
  parameters: { pseudo: { hover: true } },
  render: InputPadrao.render,
};

export const InputFoco: Story = {
  name: "Input / foco",
  parameters: { pseudo: { focusVisible: true } },
  render: InputPadrao.render,
};

export const InputDesabilitado: Story = {
  name: "Input / desabilitado",
  render: () => (
    <FormField label="E-mail">{(c) => <Input disabled defaultValue="autor@estudos.local" {...c} />}</FormField>
  ),
};

export const InputErro: Story = {
  name: "Input / erro",
  render: () => (
    <FormField label="E-mail" required error="Digite um e-mail válido, como nome@exemplo.com.">
      {(c) => <Input type="email" defaultValue="abc" {...c} />}
    </FormField>
  ),
};

export const TextareaPadrao: Story = {
  name: "Textarea",
  render: () => (
    <FormField label="Anotações" helper="Salvamos automaticamente.">
      {(c) => <Textarea rows={4} placeholder="Escreva aqui…" {...c} />}
    </FormField>
  ),
};

export const TextareaErro: Story = {
  name: "Textarea / erro",
  render: () => (
    <FormField label="Anotações" error="Não foi possível salvar.">
      {(c) => <Textarea rows={4} {...c} />}
    </FormField>
  ),
};

export const SelectPadrao: Story = {
  name: "Select",
  render: () => (
    <FormField label="Módulo">
      {(c) => (
        <Select defaultValue="sintaxe">
          <SelectTrigger {...c}>
            <SelectValue placeholder="Escolha um módulo" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="sintaxe">Sintaxe</SelectItem>
            <SelectItem value="interfaces">Interfaces</SelectItem>
          </SelectContent>
        </Select>
      )}
    </FormField>
  ),
};

export const SelectDesabilitado: Story = {
  name: "Select / desabilitado",
  render: () => (
    <FormField label="Módulo">
      {(c) => (
        <Select disabled>
          <SelectTrigger {...c}>
            <SelectValue placeholder="Indisponível" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="x">X</SelectItem>
          </SelectContent>
        </Select>
      )}
    </FormField>
  ),
};

export const CheckboxEstados: Story = {
  name: "Checkbox",
  render: () => (
    <div className="space-y-3">
      {[
        { id: "c1", label: "Desmarcado" },
        { id: "c2", label: "Marcado", defaultChecked: true },
        { id: "c3", label: "Desabilitado", disabled: true },
      ].map(({ id, label, ...rest }) => (
        <div key={id} className="flex items-center gap-2">
          <Checkbox id={id} {...rest} />
          <Label htmlFor={id}>{label}</Label>
        </div>
      ))}
    </div>
  ),
};

export const RadioEstados: Story = {
  name: "RadioGroup",
  render: () => (
    <RadioGroup defaultValue="b" aria-label="Escolha uma opção">
      {["a", "b", "c"].map((v) => (
        <div key={v} className="flex items-center gap-2">
          <RadioGroupItem value={v} id={`r-${v}`} disabled={v === "c"} />
          <Label htmlFor={`r-${v}`}>Opção {v.toUpperCase()}</Label>
        </div>
      ))}
    </RadioGroup>
  ),
};

export const SwitchEstados: Story = {
  name: "Switch",
  render: () => (
    <div className="space-y-3">
      <div className="flex items-center gap-2">
        <Switch id="s1" />
        <Label htmlFor="s1">Desligado</Label>
      </div>
      <div className="flex items-center gap-2">
        <Switch id="s2" defaultChecked />
        <Label htmlFor="s2">Ligado</Label>
      </div>
      <div className="flex items-center gap-2">
        <Switch id="s3" disabled />
        <Label htmlFor="s3">Desabilitado</Label>
      </div>
    </div>
  ),
};

export const Escuro: Story = {
  name: "Tema escuro",
  globals: { theme: "dark" },
  render: () => (
    <div className="max-w-sm space-y-4">
      <FormField label="E-mail" error="Digite um e-mail válido.">
        {(c) => <Input defaultValue="abc" {...c} />}
      </FormField>
      <FormField label="Anotações">{(c) => <Textarea rows={3} {...c} />}</FormField>
      <div className="flex items-center gap-2">
        <Checkbox id="d1" defaultChecked />
        <Label htmlFor="d1">Marcado</Label>
        <Switch aria-label="Ligado" defaultChecked />
      </div>
    </div>
  ),
};
