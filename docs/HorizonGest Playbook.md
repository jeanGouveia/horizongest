# PLAYBOOK DO HORIZONGEST

Versão 1.0

Guia Oficial de Trabalho

---

# Objetivo

Este documento explica como trabalhar dentro do HorizonGest.

O Manual Mestre explica COMO O SISTEMA FOI CONSTRUÍDO.

O Playbook explica COMO O TIME DEVE TRABALHAR.

Todo desenvolvedor ou IA deve seguir este documento.

---

# Fluxo Oficial

Toda evolução do HorizonGest segue esta ordem:

Ideia

↓

Análise

↓

Arquitetura

↓

Planejamento

↓

Implementação

↓

Testes

↓

Revisão

↓

Documentação

↓

Commit

↓

Merge

Nunca inverter esta ordem.

---

# Antes de começar qualquer tarefa

Responder:

O problema realmente existe?

Já existe solução?

Vai quebrar arquitetura?

Existe documentação?

Já existe decisão registrada?

---

# Toda tarefa começa com

1.
Entender o problema.

2.
Definir impacto.

3.
Definir arquitetura.

4.
Só então escrever código.

---

# Nunca fazer

Programar primeiro.

Pensar depois.

---

# Toda Sprint

Uma sprint deve conter:

Objetivo

Escopo

Arquivos afetados

Critérios de aceite

Critérios de teste

Impacto esperado

---

# Durante a Sprint

Sempre:

Commits pequenos.

Mudanças pequenas.

Validação constante.

Nunca alterar dezenas de arquivos sem necessidade.

---

# Ao finalizar

Executar:

Backend

go test ./...

Frontend

npm run check

Build

Testes manuais

Fluxo completo

---

# Revisão

Toda revisão deve verificar:

Arquitetura

↓

Segurança

↓

Multi-Tenant

↓

White Label

↓

Performance

↓

UX

↓

Código

---

# Commit

Nunca gigantesco.

Preferir:

feat:

fix:

refactor:

docs:

test:

style:

perf:

chore:

---

# Exemplo

feat(stock): adicionar baixa automática ao finalizar pedido

fix(order): corrigir sequência por empresa

docs(ai): atualizar contexto do projeto

---

# Pull Request

Sempre conter:

Objetivo

Problema

Solução

Arquivos alterados

Testes executados

Impacto

Rollback

---

# Quando atualizar AI_CONTEXT

Sempre que mudar:

Arquitetura

Fluxo

Módulos

Domínio

Serviços

Regras

---

# Quando atualizar DECISIONS

Sempre que existir decisão arquitetural.

Nunca registrar bugs.

Nunca registrar tarefas.

Somente decisões.

---

# Quando atualizar o Manual

Sempre que surgir uma regra permanente.

Nunca para mudanças temporárias.

---

# Fluxo com IA

A IA nunca recebe apenas:

"faz isso"

Ela recebe:

objetivo

↓

contexto

↓

restrições

↓

arquitetura

↓

critério de aceite

↓

resultado esperado

---

# Ordem oficial das IAs

ChatGPT

↓

Arquitetura

↓

Engenharia

↓

Revisão

↓

Prompt

↓

Cascade

↓

Implementação

↓

Correções

↓

Testes

---

# Fluxo recomendado

ChatGPT

↓

Prompt

↓

Cascade

↓

Implementação

↓

ChatGPT

↓

Auditoria

↓

Correção

↓

Commit

---

# Quando parar

Se um fluxo quebrar.

Não improvisar.

Corrigir.

Depois continuar.

---

# Roadmap

Sempre seguir:

Correções

↓

Fluxos

↓

UX

↓

Performance

↓

Funcionalidades

↓

Integrações

↓

IA

Nunca inverter.

---

# Filosofia

Velocidade não entrega qualidade.

Arquitetura entrega velocidade.

---

Fim do Playbook.
