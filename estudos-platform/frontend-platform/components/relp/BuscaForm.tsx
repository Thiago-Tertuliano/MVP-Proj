"use client";

import * as React from "react";
import { useRouter } from "next/navigation";

import { SearchBar } from "@/components/relp/SearchBar";

/** Campo de busca da própria página /busca (preenchido com o termo atual). */
export function BuscaForm({ termo }: { termo: string }) {
  const router = useRouter();
  return (
    <SearchBar
      defaultValue={termo}
      placeholder="O que você quer aprender?"
      label="Termo da busca"
      onSubmit={(t) => router.push(`/busca?q=${encodeURIComponent(t)}`)}
    />
  );
}
