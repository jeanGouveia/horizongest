# HORIZONGEST
## Engineering Handbook
### Constituição Técnica do Projeto

Versão: 1.0

Status:
Documento Oficial da Arquitetura

---

# Sobre este documento

Este documento representa a maior autoridade técnica do HorizonGest.

Toda decisão de engenharia deve ser compatível com este manual.

Caso exista conflito entre:

- código
- documentação antiga
- comentários
- sugestões de IA

este documento prevalece.

Ele define:

- arquitetura
- princípios
- engenharia
- padrões
- regras imutáveis
- forma correta de evolução do sistema

Este documento não explica apenas COMO fazer.

Ele explica POR QUE fazer.

Seu objetivo é preservar a qualidade da arquitetura por muitos anos.

---

# Objetivo do HorizonGest

O HorizonGest não é um CRUD.

O HorizonGest não é um sistema de telas.

O HorizonGest não é uma coleção de funcionalidades.

O HorizonGest é uma plataforma ERP multi-tenant construída para permitir que diferentes empresas administrem completamente sua operação utilizando a mesma infraestrutura.

Todo desenvolvimento deve fortalecer essa visão.

Nunca devemos adicionar funcionalidades apenas porque parecem interessantes.

Toda funcionalidade deve resolver um problema operacional real.

---

# Missão

Fornecer uma plataforma de gestão extremamente confiável.

A confiança do usuário vale mais do que qualquer funcionalidade.

Se uma decisão aumenta funcionalidades mas reduz previsibilidade, ela deve ser rejeitada.

---

# Valores Fundamentais

Toda decisão técnica deve respeitar, nesta ordem:

1. Integridade dos dados
2. Segurança
3. Arquitetura
4. Isolamento entre empresas
5. Clareza do código
6. Performance
7. Experiência do usuário
8. Estética

A ordem nunca deve ser invertida.

---

# Filosofia de Engenharia

O HorizonGest segue cinco princípios fundamentais.

## 1. Clareza acima de inteligência

Código inteligente normalmente envelhece mal.

Código explícito quase sempre sobrevive por muitos anos.

Sempre prefira:

- nomes claros
- responsabilidades pequenas
- fluxo previsível

ao invés de:

- abstrações complexas
- mágicas
- metaprogramação desnecessária

---

## 2. Simplicidade vence

A solução mais simples que atende completamente o problema deve ser escolhida.

Nunca implementar complexidade "para o futuro".

O futuro será implementado quando chegar.

---

## 3. Toda regra possui um único lugar

Uma regra de negócio nunca pode existir duplicada.

Cada regra pertence exatamente a um Service.

Nunca ao Frontend.

Nunca ao Handler.

Nunca ao Repository.

---

## 4. Backend é autoridade

O Frontend nunca decide regras.

O Frontend nunca valida negócio.

O Frontend nunca calcula processos críticos.

O Backend sempre é a fonte da verdade.

---

## 5. Evolução sem destruição

Toda evolução deve preservar:

- arquitetura
- contratos
- isolamento
- previsibilidade

O sistema deve crescer sem perder identidade.

---

# Definição de ERP

Para este projeto, ERP significa:

Sistema que controla completamente um processo operacional.

Não é suficiente cadastrar dados.

É necessário representar regras reais.

Exemplo:

Pedido não é apenas uma tabela.

Pedido representa:

- estoque
- produção
- financeiro
- histórico
- auditoria
- métricas

---

# Escopo

O HorizonGest foi desenhado para crescer continuamente.

Ele deverá suportar dezenas de módulos.

Exemplos:

- Financeiro

- Compras

- Estoque

- Produção

- Delivery

- CRM

- RH

- Agenda

- Marketing

- Fiscal

- Dashboard

- BI

Todos esses módulos devem compartilhar a mesma arquitetura.

Nenhum módulo pode criar sua própria arquitetura.

---

# O que significa arquitetura

Arquitetura não é organização de pastas.

Arquitetura é o conjunto de regras que impede o sistema de se degradar.

Quando duas pessoas diferentes escrevem código seguindo as mesmas regras, o sistema continua parecendo escrito por uma única equipe.

Esse é o objetivo.

---

# O que nunca deve acontecer

Jamais devemos permitir:

- regras espalhadas

- código duplicado

- acoplamento crescente

- dependências circulares

- quebra do isolamento entre empresas

- decisões locais que prejudiquem o projeto inteiro

---

# Definição de qualidade

Qualidade não significa menos linhas.

Qualidade não significa código moderno.

Qualidade significa:

- previsibilidade

- facilidade de manutenção

- facilidade de testes

- facilidade de evolução

- baixo risco

---

# Definição de dívida técnica

Dívida técnica é toda decisão que:

- reduz previsibilidade

- aumenta acoplamento

- duplica responsabilidade

- dificulta manutenção

Mesmo que funcione.

---

# Regra de Ouro

Nenhuma decisão deve ser tomada pensando apenas na tarefa atual.

Toda decisão deve responder:

"Isso ainda fará sentido daqui a cinco anos?"

Se a resposta for "não sei",

a decisão provavelmente está errada.

---

# Visão de Longo Prazo

O HorizonGest foi projetado para sobreviver durante muitos anos.

Isso exige:

- arquitetura estável

- documentação forte

- padrões claros

- evolução incremental

Nunca reescrever por moda.

Nunca trocar tecnologia sem necessidade.

Nunca quebrar contratos existentes apenas porque existe uma forma "mais bonita".

---

Fim da Parte 1

# PARTE 2
# Constituição da Arquitetura

---

# Capítulo 1
# A Constituição da Arquitetura

Toda arquitetura possui princípios.

O HorizonGest possui leis.

Essas leis não são sugestões.

São restrições obrigatórias.

Sempre que uma alteração for proposta, ela deve ser comparada contra esta Constituição.

Se existir conflito, a alteração deve ser descartada.

Mesmo que ela funcione.

Mesmo que reduza linhas de código.

Mesmo que aumente performance.

Arquitetura sempre possui prioridade maior.

---

# Artigo 1
## O domínio pertence ao Backend

Todo comportamento do sistema pertence ao Backend.

O Frontend apenas representa esse comportamento.

Nunca o contrário.

Exemplos de domínio:

- criar pedido

- cancelar pedido

- validar estoque

- calcular lucro

- calcular impostos

- fechar caixa

- validar permissões

- gerar relatórios

Nada disso pertence ao Frontend.

---

# Artigo 2
## O Frontend representa estados

O Frontend possui apenas quatro responsabilidades.

Apresentação.

Experiência do usuário.

Navegação.

Estado visual.

Nada além disso.

Sempre que surgir uma dúvida:

"Essa lógica pertence ao Frontend?"

A resposta deve ser:

Provavelmente não.

---

# Artigo 3
## O Backend nunca conhece interface

O Backend nunca deve saber:

- qual framework frontend existe

- qual biblioteca existe

- qual tela chamou

- qual botão foi clicado

- qual componente disparou a requisição

Para o Backend existe apenas:

Request

Response

Domínio

Persistência

---

# Artigo 4
## A arquitetura é direcional

Toda dependência possui direção única.

A direção nunca muda.

Nunca.

Fluxo permitido:

Frontend

↓

HTTP

↓

Handler

↓

Service

↓

Repository

↓

Banco

O fluxo contrário nunca é permitido.

---

# Artigo 5
## Nenhuma camada pode pular outra

É proibido:

Frontend chamar Repository.

Frontend chamar banco.

Handler chamar banco.

Handler chamar Repository.

Repository chamar Service.

Repository chamar Handler.

Service chamar Handler.

Banco chamar qualquer camada.

Toda camada conversa apenas com sua vizinha imediata.

---

# Artigo 6
## Handler é adaptador

O Handler não possui inteligência.

Ele apenas adapta HTTP para domínio.

Responsabilidades permitidas:

ler parâmetros

ler headers

ler body

validar formato

converter tipos

chamar Service

converter resposta

retornar HTTP

Nada além disso.

---

# Artigo 7
## Service é o cérebro

Toda inteligência mora aqui.

Toda regra.

Toda decisão.

Toda validação.

Todo cálculo.

Toda política.

Toda autorização de negócio.

Todo fluxo operacional.

Sempre que surgir uma regra nova:

ela nasce em um Service.

---

# Artigo 8
## Repository não toma decisões

Repository apenas executa operações de persistência.

Ele não interpreta regras.

Não calcula.

Não decide.

Não valida.

Não conhece UX.

Não conhece negócio.

Ele apenas responde perguntas como:

Buscar.

Salvar.

Atualizar.

Remover.

Consultar.

---

# Artigo 9
## O Banco é armazenamento

O banco não contém regras.

Constraints são permitidas.

Índices são permitidos.

Foreign Keys são obrigatórias.

Triggers devem ser evitadas.

Stored Procedures devem ser evitadas.

Toda inteligência permanece na aplicação.

---

# Artigo 10
## Nenhuma regra pode ser duplicada

Uma regra existe exatamente uma vez.

Se ela aparece em dois lugares,

um deles está errado.

Duplicação é um defeito arquitetural.

---

# Artigo 11
## Toda regra possui dono

Exemplo.

Preço do produto.

Pertence ao ProductService.

Nunca ao OrderService.

Nunca ao Frontend.

Nunca ao Repository.

Cada regra possui proprietário.

---

# Artigo 12
## Services representam Domínios

Um Service não representa tabelas.

Representa responsabilidades.

Exemplo correto:

OrderService

InventoryService

FinanceService

NotificationService

Exemplo incorreto:

OrdersTableService

ProductsTableService

---

# Artigo 13
## Um Service pequeno vale mais que um Service inteligente

Quando um Service cresce demais,

ele perde identidade.

Preferimos:

cinco Services claros

do que

um Service gigante.

---

# Artigo 14
## O sistema deve parecer escrito por uma única pessoa

Mesmo que cinquenta desenvolvedores participem,

o código deve possuir identidade única.

Mesmo padrão.

Mesmo estilo.

Mesmo fluxo.

Mesmo nível de abstração.

---

# Artigo 15
## Comentários explicam decisões

Comentários nunca devem explicar código.

Código já explica código.

Comentários devem explicar:

por que

e não

como.

Exemplo ruim:

// soma total

total += valor

Exemplo bom:

// O valor é acumulado antes do desconto
// porque o imposto incide sobre o subtotal.

---

# Artigo 16
## Clareza vence elegância

Nunca escrevemos código para impressionar.

Escrevemos código para sobreviver.

O melhor código é aquele que outro desenvolvedor entende em cinco minutos.

Não aquele que parece genial.

---

# Artigo 17
## Toda decisão deve reduzir risco

Toda alteração deve responder:

Essa mudança aumenta ou reduz risco?

Se aumentar,

ela provavelmente está errada.

---

# Artigo 18
## Arquitetura não negocia

Performance negocia.

UX negocia.

Visual negocia.

Arquitetura não.

Arquitetura é permanente.

---

# Resumo deste capítulo

Se um desenvolvedor esquecer todo o restante do manual,

mas lembrar destes artigos,

ele ainda conseguirá manter o HorizonGest saudável por muitos anos.

Fim da Parte 2

# PARTE 3
# Anatomia do HorizonGest

---

# Capítulo 2
# Anatomia da Arquitetura

Toda arquitetura saudável possui responsabilidades extremamente claras.

Quando uma responsabilidade não possui dono,
ela acaba aparecendo em vários lugares.

É assim que nasce a dívida técnica.

O HorizonGest elimina esse problema definindo exatamente
qual camada possui qual responsabilidade.

Cada camada possui:

- objetivo
- responsabilidades
- permissões
- proibições

---

# Visão Geral

Fluxo oficial:

Frontend

↓

HTTP

↓

Middleware

↓

Handler

↓

Service

↓

Repository

↓

Database

Nada pode alterar essa sequência.

---

# Frontend

## Objetivo

Representar visualmente o estado do sistema.

Nada além disso.

---

## Responsabilidades

Mostrar informações.

Capturar ações do usuário.

Enviar requisições.

Exibir respostas.

Gerenciar estado visual.

Gerenciar componentes.

Controlar navegação.

Experiência do usuário.

Responsividade.

Acessibilidade.

---

## O Frontend PODE

Mostrar loading.

Mostrar erro.

Mostrar sucesso.

Controlar modais.

Controlar menus.

Controlar rotas.

Validar formato de campos.

Mascarar inputs.

Controlar animações.

Controlar foco.

---

## O Frontend NÃO PODE

Calcular impostos.

Calcular lucro.

Calcular estoque.

Definir permissões.

Decidir regras.

Criar IDs.

Gerar números de pedidos.

Validar regras de negócio.

Calcular descontos.

Calcular comissão.

Persistir diretamente.

Executar SQL.

Interpretar domínio.

---

## Exemplo correto

Botão:

Salvar

↓

POST /products

↓

Backend responde

↓

Atualiza interface.

---

## Exemplo incorreto

Botão:

Salvar

↓

Calcula margem

↓

Calcula comissão

↓

Decide desconto

↓

Envia resultado

Toda essa lógica pertence ao Backend.

---

# Middleware

## Objetivo

Interceptar requisições.

Nunca implementar negócio.

---

## Responsabilidades

Autenticação.

Autorização técnica.

Rate Limit.

Logs.

Recovery.

CORS.

Headers.

Contexto.

Tenant.

---

## Middleware NÃO PODE

Consultar Repository.

Executar regra de negócio.

Alterar domínio.

Tomar decisões comerciais.

---

# Handler

## Objetivo

Traduzir HTTP para domínio.

É apenas um adaptador.

---

## Responsabilidades

Ler parâmetros.

Ler QueryString.

Ler Headers.

Ler JSON.

Validar formato.

Converter tipos.

Chamar Service.

Converter Response.

Retornar HTTP Status.

---

## Handler NÃO PODE

Consultar banco.

Executar SQL.

Calcular estoque.

Calcular financeiro.

Criar pedidos.

Cancelar pedidos.

Atualizar domínio.

Tomar decisões.

---

## Handler ideal

Recebe.

Valida formato.

Entrega.

Responde.

Fim.

---

# Service

## Objetivo

Executar domínio.

Toda inteligência mora aqui.

---

## Responsabilidades

Validação de negócio.

Fluxos.

Regras.

Políticas.

Permissões.

Cálculos.

Integrações.

Eventos.

Orquestração.

---

## O Service PODE

Consultar vários repositories.

Executar transações.

Controlar rollback.

Chamar serviços externos.

Executar cálculos.

Controlar regras.

Gerar eventos.

---

## O Service NÃO PODE

Executar SQL.

Conhecer HTTP.

Conhecer HTML.

Conhecer Svelte.

Conhecer componentes.

Conhecer telas.

Conhecer CSS.

---

# Repository

## Objetivo

Persistência.

Nada além.

---

## Responsabilidades

Buscar.

Salvar.

Atualizar.

Excluir.

Consultar.

Paginar.

Filtrar.

---

## Repository NÃO PODE

Calcular.

Validar.

Autorizar.

Notificar.

Tomar decisões.

Conhecer Frontend.

Conhecer HTTP.

---

# Domain

## Objetivo

Representar entidades.

Não comportamento.

---

## Domain contém

Structs.

Enums.

Constantes.

Interfaces.

Tipos.

---

## Domain NÃO contém

Queries.

HTTP.

JSON.

Banco.

Framework.

---

# Database

## Objetivo

Persistir informações.

---

## Banco deve conter

Constraints.

Foreign Keys.

Indexes.

Unique Keys.

Not Null.

Check Constraints.

---

## Banco NÃO deve conter

Regra comercial.

Fluxo.

Triggers complexas.

Stored Procedures de negócio.

---

# Configuração

Existe diferença entre:

Configuração

e

Código.

---

Configuração muda.

Código não.

---

Exemplos de configuração

Brand.

SMTP.

Rate Limit.

Timeout.

Logo.

Nome da Plataforma.

Tema.

Feature Flags.

---

Exemplos de código

Cálculo de estoque.

Pedido.

Financeiro.

Fluxo de caixa.

RBAC.

Nunca devem virar configuração.

---

# Migrations

Migrations são históricas.

Jamais editar uma migration aplicada.

Nunca.

Nova alteração.

Nova migration.

Sempre.

---

# DTOs

DTO existe apenas para transportar dados.

Nunca deve possuir regra.

Nunca deve possuir comportamento.

---

# Events

Eventos representam fatos.

Nunca comandos.

Correto:

PedidoCriado

CompraFinalizada

ProdutoExcluído

PagamentoConfirmado

Errado:

CriarPedido

ExcluirProduto

AtualizarCliente

---

# Controllers

No HorizonGest usamos Handlers.

Não Controllers.

O motivo é simples.

Controller sugere lógica.

Handler sugere adaptação.

Essa escolha é intencional.

---

# Dependências

Toda camada conhece apenas quem está abaixo dela.

Frontend

↓

Handler

↓

Service

↓

Repository

↓

Database

Nunca o contrário.

---

# Acoplamento permitido

Handler

↓

Service

Service

↓

Repository

Repository

↓

Database

---

# Acoplamento proibido

Frontend

↓

Repository

Repository

↓

Service

Database

↓

Service

Service

↓

Handler

---

# Responsabilidade Única

Se um arquivo responde por duas perguntas diferentes,

ele provavelmente deveria ser dividido.

---

# Regra prática

Quando surgir uma dúvida sobre onde escrever um código,
faça apenas uma pergunta:

"Isso representa domínio ou infraestrutura?"

Se for domínio,

vai para Service.

Se for infraestrutura,

vai para Repository ou Middleware.

Se for interface,

vai para Frontend.

Essa regra resolve aproximadamente 90% das dúvidas arquiteturais.

Fim da Parte 3

# PARTE 4
# Arquitetura Multi-Tenant

---

# Capítulo 3
# A Arquitetura Multi-Tenant do HorizonGest

O HorizonGest foi concebido desde o primeiro dia como um SaaS multi-tenant.

Essa decisão influencia absolutamente toda a arquitetura.

Nada no sistema pode ser desenvolvido ignorando este conceito.

Todo código escrito deve assumir que:

- existirão milhares de empresas
- todas compartilharão a mesma aplicação
- todas compartilharão o mesmo banco
- nenhuma empresa poderá descobrir a existência das demais

Essa é uma regra absoluta.

---

# Conceitos Fundamentais

Existem apenas dois níveis de operação.

Platform.

Tenant.

Nunca criar um terceiro nível.

---

# O que é Platform

Platform representa o proprietário do sistema.

A Platform administra o HorizonGest.

Ela não administra restaurantes.

Ela administra clientes.

Exemplos:

- cadastro de empresas

- planos

- cobrança

- branding

- configurações globais

- auditoria global

- usuários da plataforma

- monitoramento

- feature flags

Platform nunca vende produtos.

Nunca faz pedidos.

Nunca controla estoque.

---

# O que é Tenant

Tenant representa uma empresa cliente.

Cada empresa vive isoladamente.

Cada empresa possui:

- usuários

- pedidos

- estoque

- produtos

- clientes

- fornecedores

- financeiro

- compras

- produção

Tudo pertence exclusivamente ao Tenant.

---

# Company

No HorizonGest, Tenant é representado por Company.

Toda Company é completamente independente.

Mesmo que duas empresas tenham:

- mesmo nome

- mesmo CNPJ

- mesmos usuários

- mesmos produtos

elas continuam sendo entidades independentes.

---

# CompanyID

CompanyID é o principal elemento arquitetural do sistema.

Ele representa o isolamento de dados.

Sem CompanyID não existe multi-tenancy.

---

# Regra Fundamental

Toda entidade pertencente ao Tenant deve possuir CompanyID.

Sem exceções.

---

# Exemplos

Possuem CompanyID

Products

Orders

Ingredients

Categories

Purchases

StockMovements

Finance

Customers

Suppliers

Employees

Permissions

Notifications

Media

BusinessProfile

Theme

Tudo isso pertence ao Tenant.

---

# Entidades Globais

Algumas entidades pertencem apenas à Platform.

Essas não possuem CompanyID.

Exemplos:

PlatformUser

PlatformSession

PlatformAudit

Plan

PlatformBrand

GlobalConfig

FeatureFlag

SystemVersion

---

# Nunca misturar

Jamais criar uma entidade híbrida.

Ou ela pertence à Platform.

Ou pertence ao Tenant.

Nunca aos dois.

---

# Company Context

Toda requisição autenticada deve possuir contexto de Company.

Exemplo:

JWT

↓

Middleware

↓

CompanyID

↓

Context

↓

Handler

↓

Service

↓

Repository

↓

Banco

O CompanyID deve viajar durante toda a requisição.

Nunca deve ser recalculado.

Nunca deve ser recebido do Frontend.

---

# O Frontend nunca escolhe CompanyID

Essa talvez seja a regra mais importante do sistema.

O Frontend nunca envia CompanyID.

Nunca.

Jamais.

O Backend descobre automaticamente qual empresa está autenticada.

Isso impede:

- fraude

- troca manual

- acesso indevido

- manipulação de requisições

---

# Isolamento de Dados

Todo Repository deve aplicar CompanyID.

Sempre.

Mesmo quando parecer desnecessário.

---

# Correto

SELECT *

FROM products

WHERE company_id = ?

---

# Errado

SELECT *

FROM products

Esse tipo de erro é considerado crítico.

---

# Busca por ID

Mesmo quando buscando por ID.

Correto:

SELECT *

FROM products

WHERE id = ?

AND company_id = ?

Nunca apenas:

WHERE id = ?

---

# Atualizações

Toda atualização deve respeitar CompanyID.

UPDATE

↓

id

+

company_id

Nunca atualizar apenas pelo ID.

---

# Exclusões

Toda exclusão deve respeitar CompanyID.

DELETE

↓

id

+

company_id

---

# Numeração Interna

Existem dois tipos de identificadores.

---

# Internal ID

Representa a chave primária.

É global.

Nunca é mostrado ao usuário.

Serve apenas para relacionamentos internos.

---

# Business Number

Representa números visíveis.

Exemplos:

Pedido 15

Compra 8

Venda 230

Esses números pertencem ao Tenant.

Nunca são globais.

---

# Exemplo

Empresa A

Pedido 1

Pedido 2

Pedido 3

Empresa B

Pedido 1

Pedido 2

Empresa C

Pedido 1

Todos corretos.

---

# IDs Globais

IDs globais nunca podem aparecer em:

Interface.

PDF.

Relatórios.

Notas.

Comprovantes.

URLs públicas.

Mensagens.

Sempre utilizar Business Number.

---

# Segurança

Uma empresa nunca pode descobrir:

quantos pedidos existem em outra empresa.

quantos produtos existem.

quantos usuários existem.

quais planos estão ativos.

quais empresas utilizam o sistema.

Toda informação deve permanecer isolada.

---

# Auditoria

Logs globais pertencem à Platform.

Logs operacionais pertencem ao Tenant.

Nunca misturar.

---

# Cache

Quando existir cache,

ele também deve respeitar CompanyID.

Exemplo correto:

product:company:15:id:80

Nunca:

product:80

---

# Eventos

Eventos também pertencem ao Tenant.

PedidoCriado

↓

CompanyID

↓

Evento

Nunca publicar eventos sem contexto.

---

# Exportações

Toda exportação deve ser filtrada por Company.

---

# Session Management

O HorizonGest possui uma arquitetura profissional de gerenciamento de sessões que garante segurança, consistência e previsibilidade em todo o ciclo de vida da autenticação.

## Tipos de Sessão

Existem dois tipos de sessão:

### Platform Session

Responsável pelo acesso à Plataforma.

**Características:**
- Exige login
- Possui expiração
- Nunca pode ser restaurada sem validação do backend
- Perde validade após reinício do backend
- Perde validade após logout
- Perde validade após expiração do JWT

**Armazenamento:**
- Cookie: `platform_auth_token`

### Tenant Session

Responsável pelo acesso à empresa.

**Características:**
- Só existe enquanto existir uma Platform Session válida
- Nunca pode sobreviver sozinha
- Nunca pode sobreviver ao logout
- Nunca pode sobreviver ao reinício do backend
- Nunca pode sobreviver à troca de empresa

**Armazenamento:**
- Cookie: `auth_token`
- LocalStorage: `impersonation`

## Managers

### SessionManager

**Arquivo:** `frontend/src/lib/managers/sessionManager.ts`

**Responsabilidades:**
- Validar sessão na inicialização
- Gerenciar Platform Session
- Gerenciar Tenant Session (via TenantSessionManager)
- Executar logout completo
- Destruir todas as sessões

### TenantSessionManager

**Arquivo:** `frontend/src/lib/managers/tenantSessionManager.ts`

**Responsabilidades:**
- Gerenciar entrada em empresa (`enterCompany`)
- Gerenciar saída de empresa (`leaveCompany`)
- Gerenciar troca de empresa (`switchCompany`)
- Destruir contexto completo (`destroy`)

## Regras Fundamentais

### 1. Único Ponto de Troca
Nenhum componente deve trocar empresa diretamente. Toda troca deve passar pelo `TenantSessionManager`.

### 2. Navegação Apenas Após Contexto Consistente
Nunca navegar antes do contexto estar consistente. Os managers garantem que:
- Contexto anterior foi destruído
- Novo token foi obtido
- Stores foram limpas
- Caches foram limpos
- Novo contexto foi hidratado
- Branding foi carregado
- Permissões foram carregadas
- Empresa foi carregada

### 3. Validação de Sessão
Nunca abrir Dashboard apenas porque existe token salvo. Sempre validar sessão no backend na inicialização.

### 4. Tratamento de 401
Qualquer endpoint que retorna 401 deve automaticamente destruir todas as sessões e redirecionar para login.

### 5. Prevenção de Sessões Órfãs
Backend deve automaticamente limpar sessões stale (mais de 24 horas) antes de criar novas impersonations.

### 6. Storage Keys Centralizadas
Todas as chaves de storage (cookies, localStorage, sessionStorage) devem ser centralizadas em `frontend/src/lib/constants/storage-keys.ts`.

**Rationale:**
- Previne typos em strings de storage
- Facilita refatoração quando chaves precisam mudar
- Documenta o propósito de cada chave
- Fornece type safety através de constantes TypeScript

**Regra:** NUNCA usar strings literais para storage keys. Sempre importar do arquivo de constantes.

### 7. Limpeza Granular de Cache
Nunca usar `sessionStorage.clear()` ou `localStorage.clear()` globalmente. Limpar apenas chaves específicas do contexto.

**Rationale:**
- `sessionStorage.clear()` remove TODOS os dados de sessão, incluindo dados da Platform
- Limpeza granular garante que apenas dados do Tenant são removidos
- Previne remoção acidental de dados da Platform Session
- Mais previsível e maintenível

**Regra:** Implementar métodos específicos como `clearTenantSessionStorage()` que removem apenas chaves do Tenant.

### 8. Compatibilidade de Navegador
Usar apenas APIs suportadas por todos os navegadores alvo (Chrome, Firefox, Safari, Edge).

**Rationale:**
- APIs experimentais podem não funcionar em todos os navegadores
- APIs experimentais podem mudar ou ser removidas sem aviso
- Usar APIs padrão garante compatibilidade cross-browser
- Reduz risco de erros em produção

**Regra:** Verificar suporte em caniuse.com antes de usar novas APIs. Evitar APIs marcadas como "experimental".

### 9. Tratamento Diferenciado de Erros
Separar erros por tipo: infrastructure, session, backend, ui. Fornecer mensagens específicas para cada tipo.

**Rationale:**
- Diferentes tipos de erro requerem diferentes mensagens ao usuário
- Erros de infraestrutura (rede) precisam de tratamento diferente de erros de backend
- Erros de validação de sessão devem ser tratados diferentemente de erros de UI
- Melhora debugging e logging

**Regra:** Usar classes de erro apropriadas para diferentes cenários. Fornecer mensagens amigáveis ao usuário baseadas no tipo de erro.

### 10. Sem @ts-ignore
Nunca usar `@ts-ignore` para silenciar erros TypeScript.

**Rationale:**
- `@ts-ignore` silencia erros TypeScript sem corrigir a causa raiz
- Esconde bugs potenciais e problemas de type safety
- Torna o codebase menos mantenível
- Impede TypeScript de fazer seu trabalho

**Regra:** Corrigir erros TypeScript adequadamente em vez de usar `@ts-ignore`.

## Documentação Completa

Para documentação detalhada sobre Session Management, consulte:
- `docs/05-development/SESSION_MANAGEMENT.md`
- `docs/05-development/SESSION_TESTING.md`

---

Jamais exportar dados globais para Tenant.

---

# Importações

Toda importação deve associar automaticamente CompanyID.

Nunca confiar em CompanyID vindo do arquivo.

---

# Backup

Backups podem existir em dois níveis.

Platform Backup.

Tenant Backup.

Nunca misturar.

---

# Testes

Todo teste envolvendo Tenant deve verificar:

Empresa A não vê Empresa B.

Empresa B não altera Empresa A.

Empresa C não consulta Empresa A.

Esse teste nunca deve ser removido.

---

# Regra Suprema

Se existir qualquer dúvida entre:

Facilidade de implementação

e

Isolamento de dados

sempre escolher isolamento.

Mesmo que o código fique maior.

Mesmo que a performance diminua.

Mesmo que seja necessário escrever mais SQL.

O isolamento entre empresas é um dos pilares permanentes do HorizonGest.

Fim da Parte 4


# PARTE 5
# Organização Física do Projeto

---

# Capítulo 4
# Organização Física

Uma arquitetura limpa começa pela organização do código.

A estrutura de diretórios do HorizonGest não existe por estética.

Ela existe para reduzir a carga cognitiva.

Qualquer desenvolvedor deve conseguir localizar um arquivo em poucos segundos.

Se alguém precisa procurar muito para encontrar uma responsabilidade,
a estrutura está errada.

---

# Princípio Fundamental

Cada diretório possui uma única responsabilidade.

Nunca criar diretórios genéricos como:

utils2

helpers

misc

temp

novo

coisas

Esses nomes escondem responsabilidades.

---

# Estrutura Geral

HorizonGest

```
backend/
frontend/
docs/
scripts/
```

Nada além disso deve existir na raiz, exceto arquivos de configuração.

---

# Backend

O backend representa o domínio do sistema.

Toda inteligência mora aqui.

Estrutura:

```
backend/

cmd/

internal/

migrations/

tests/

```

---

# cmd/

Contém apenas inicialização da aplicação.

Exemplo:

```
cmd/server
```

Responsabilidades:

- iniciar servidor

- carregar configuração

- registrar dependências

- iniciar HTTP

Nunca colocar regra de negócio.

---

# internal/

Representa toda a aplicação.

Tudo que faz o sistema funcionar mora aqui.

Estrutura oficial:

```
internal/

domain/

handler/

service/

ports/

infra/

middleware/

util/

```

---

# domain/

Representa o negócio.

Contém:

Entities

Enums

Interfaces

Value Objects

Constantes

Nunca contém:

SQL

HTTP

Framework

GORM

JSON

---

# handler/

Responsável por HTTP.

Contém:

Request

Response

DTO

Conversões

Nunca:

SQL

Regra

Cálculos

---

# service/

O coração do sistema.

Toda regra nasce aqui.

Cada Service representa um domínio.

Exemplo:

```
OrderService

ProductService

FinanceService

StockService
```

Nunca criar:

```
UtilsService

CommonService

GenericService

BaseService
```

Esses nomes escondem responsabilidades.

---

# ports/

Define contratos.

Nunca implementação.

Exemplo:

```
UserRepository

OrderRepository

PaymentGateway
```

Nunca colocar lógica.

---

# infra/

Representa infraestrutura.

Estrutura:

```
infra/

repository/

database/

cache/

email/

storage/

external/

```

---

# repository/

Contém implementações GORM.

Exemplo:

```
gorm_order_repository.go

gorm_product_repository.go
```

Nunca colocar regras.

---

# database/

Responsável por conexão.

Migration.

Seeds.

Configuração.

Nunca regras.

---

# middleware/

Interceptação HTTP.

Exemplos:

JWT

RBAC

Tenant

Recovery

CORS

RateLimit

Nunca negócio.

---

# util/

Utilitários puros.

Funções matemáticas.

Conversões.

Datas.

Nunca domínio.

Nunca regras.

---

# Migrations

Toda migration deve possuir número sequencial.

Exemplo:

```
00028_add_discount_table.sql
```

Nunca:

```
nova.sql

migration.sql

alteracao.sql
```

---

# Frontend

O Frontend representa apenas interface.

Estrutura oficial:

```
src/

lib/

routes/

stores/

theme/

components/

```

---

# routes/

Representa páginas.

Nada além.

Nunca colocar componentes reutilizáveis.

---

# lib/

Biblioteca compartilhada.

Subestrutura:

```
components/

stores/

api/

utils/

theme/

types/

```

---

# components/

Componentes reutilizáveis.

Nunca páginas completas.

Nunca regras.

---

# stores/

Estado global.

Nunca regras de negócio.

Pode conter:

Sessão.

Tema.

Brand.

Idioma.

Notificações.

---

# api/

Cliente HTTP.

Nunca regras.

Nunca cálculo.

---

# theme/

Todo Design System.

Cor.

Fonte.

Espaçamento.

Shadow.

Radius.

Transition.

Nunca componentes.

---

# types/

Interfaces TypeScript.

Nada além.

---

# docs/

Toda documentação oficial.

Nunca espalhar documentação na raiz.

Estrutura oficial:

```
docs/

01-overview/

02-backend/

03-frontend/

04-platform/

05-development/

06-reference/

manuals/

```

---

# scripts/

Scripts auxiliares.

Exemplos:

Backup.

Deploy.

Seed.

Conversões.

Nunca regras.

---

# Nomeação

Arquivos seguem padrão consistente.

Correto:

```
order_service.go

product_repository.go

company_handler.go
```

Errado:

```
OrderService.go

pedidoNovo.go

teste2.go
```

---

# Organização por responsabilidade

Nunca organizar por tipo técnico.

Sempre por responsabilidade.

Correto:

```
ProductService

ProductRepository

ProductHandler
```

Todos pertencem ao domínio Product.

---

# Arquivos Grandes

Quando um arquivo ultrapassar aproximadamente 500 linhas,
ele deve ser avaliado.

Pergunta:

Ele continua representando uma única responsabilidade?

Se não,

deve ser dividido.

---

# Imports

Sempre utilizar imports absolutos do módulo.

Nunca caminhos improvisados.

Nunca duplicar dependências.

---

# Dependências

Todo arquivo deve depender apenas do necessário.

Imports não utilizados representam dívida técnica.

---

# Comentários

Comentários pertencem ao "porquê".

Nunca ao "como".

---

# Organização é Arquitetura

Mover um arquivo para uma pasta errada
também é uma quebra arquitetural.

A organização física do projeto faz parte da arquitetura.

Ela deve ser preservada durante toda a vida do HorizonGest.

---

# Regra Final

Se surgir dúvida sobre onde criar um arquivo,
a resposta nunca deve ser:

"qualquer lugar".

Existe exatamente um lugar correto para cada responsabilidade.

Encontrar esse lugar é obrigação do desenvolvedor.

Fim da Parte 5


# PARTE 6
# Convenções de Código e Engenharia

---

# Capítulo 5
# O Jeito HorizonGest de Programar

Arquitetura organiza o sistema.

Convenções organizam os desenvolvedores.

Dois programadores diferentes devem produzir código praticamente indistinguível.

O objetivo não é limitar criatividade.

O objetivo é reduzir atrito.

---

# Princípio Fundamental

Código é lido muito mais vezes do que escrito.

Escrevemos para quem vai manter.

Não para quem escreveu.

---

# Legibilidade

Sempre priorizar:

clareza

sobre

brevidade.

---

# Regra

Se um código é menor mas fica mais difícil de entender,

ele está pior.

---

# Nomes

Todo nome deve responder claramente:

"O que isto representa?"

---

# Variáveis

Correto

```
totalPrice

discountValue

customerName

employeeID

companyID
```

Errado

```
x

tmp

aux

valor2

teste

obj

itemNovo2
```

---

# Funções

Toda função deve possuir um verbo.

Exemplo

```
CreateOrder

UpdateStock

CalculateTotal

GenerateInvoice

SendEmail
```

Nunca

```
Order

Stock

Invoice

Email
```

---

# Funções pequenas

Idealmente

20 linhas.

Aceitável

50 linhas.

Acima disso

avaliar divisão.

---

# Responsabilidade Única

Uma função deve responder apenas uma pergunta.

Nunca fazer:

- validar

- salvar

- enviar email

- gerar relatório

na mesma função.

---

# Return antecipado

Preferimos:

```
if err != nil {
    return err
}
```

ao invés de

```
if err == nil {

    ...

}
```

Reduz indentação.

Melhora leitura.

---

# Else

Evitar else quando houver return.

Correto

```
if err != nil {
    return err
}

return process()
```

Errado

```
if err != nil {

    return err

} else {

    return process()

}
```

---

# Booleanos

Boolean deve responder pergunta.

Correto

```
isAdmin

isOwner

hasPermission

canDelete
```

Errado

```
admin

owner

permission

delete
```

---

# Strings mágicas

Nunca espalhar strings.

Correto

```
const OrderStatusPaid
```

Errado

```
"paid"
```

em dezenas de arquivos.

---

# Numbers mágicos

Nunca utilizar

```
7

30

100

999
```

sem contexto.

Sempre transformar em constante.

---

# Constantes

Devem representar conceitos.

Nunca apenas números.

---

# Comentários

Comentários existem para explicar decisões.

Nunca código.

---

# Exemplo ruim

```
Incrementa contador
contador++
```

---

# Exemplo correto

```
Incrementamos antes da persistência
porque o ERP exige numeração contínua.
```

---

# TODO

TODO deve possuir motivo.

Correto

```
TODO:

Migrar para Redis quando houver múltiplas instâncias.
```

Errado

```
TODO melhorar
```

---

# FIXME

Somente para defeitos conhecidos.

Nunca para melhorias.

---

# Logs

Logs representam eventos.

Nunca debug permanente.

---

# Log correto

```
Pedido criado

Company 15

Order 84

User 7
```

---

# Log ruim

```
Entrou aqui

teste

oi

123

debug
```

---

# Tratamento de erros

Nunca ignorar erro.

Errado

```
result, _ := process()
```

---

Sempre

```
result, err := process()

if err != nil {

    return err

}
```

---

# Wrapping de erros

Sempre preservar contexto.

Correto

```
return fmt.Errorf(
"CreateOrder: %w",
err,
)
```

Nunca

```
return err
```

sozinho.

---

# Panic

Nunca utilizar panic para fluxo normal.

Panic apenas para situações irrecuperáveis.

---

# Recover

Recover pertence apenas ao middleware.

Nunca ao domínio.

---

# Transações

Toda transação deve possuir início e fim claros.

Dentro da transação apenas operações que realmente precisam ser atômicas.

Nunca colocar chamadas HTTP.

Nunca enviar emails.

Nunca executar tarefas demoradas.

---

# SQL

Todo SQL deve ser previsível.

Nunca construir SQL concatenando strings.

Sempre parametrizado.

---

# Performance

Nunca otimizar antes da necessidade.

Primeiro

clareza.

Depois

correção.

Depois

performance.

---

# Cache

Cache nunca pode alterar comportamento.

Se o cache desaparecer,

o sistema deve continuar correto.

---

# Testes

Todo teste deve possuir nome descritivo.

Correto

```
TestCreateOrderWithoutStockReturnsError
```

Errado

```
Test1

TestNovo

TestPedido
```

---

# Estrutura dos testes

Arrange

Act

Assert

Sempre nessa ordem.

---

# Commits

Todo commit representa uma única intenção.

Nunca misturar:

bug

refactor

feature

documentação

na mesma mensagem.

---

# Commits seguem Conventional Commits

Exemplos

```
feat:

fix:

refactor:

docs:

test:

perf:

build:

chore:
```

---

# Pull Request

Toda PR deve responder:

O que mudou?

Por que mudou?

Qual risco?

Como testar?

---

# Revisão de Código

Durante uma revisão perguntar sempre:

Quebrou arquitetura?

Duplicou regra?

Criou acoplamento?

Existe responsabilidade errada?

Existe código morto?

Existe simplificação possível?

---

# Regra dos Cinco Anos

Antes de salvar um código,

pergunte:

"Alguém entenderá isto daqui cinco anos?"

Se a resposta for não,

reescreva.

---

# Filosofia Final

Escrevemos código para durar.

Não escrevemos código para impressionar.

O melhor elogio que um código pode receber é:

"Foi fácil entender."

Esse é o padrão HorizonGest.

Fim da Parte 6


# PARTE 7
# Conhecendo o HorizonGest

---

# Capítulo 6
# Visão Geral do Sistema

O HorizonGest é uma plataforma SaaS multiempresa.

Seu objetivo é permitir que diferentes empresas utilizem o mesmo sistema, mantendo isolamento completo de dados.

Cada empresa enxerga apenas seus próprios dados.

Toda a arquitetura foi construída em torno desse princípio.

---

# Objetivos

O HorizonGest deve ser:

- Multiempresa

- Escalável

- White Label

- Modular

- Seguro

- Offline Friendly (quando necessário)

- Fácil de manter

- Fácil de evoluir

---

# Estrutura Geral

Existem dois grandes ambientes:

Platform

e

Company.

---

# Platform

Platform administra toda a plataforma.

Responsável por:

- empresas

- planos

- usuários da plataforma

- branding

- configurações globais

- auditoria

- feature flags

- métricas gerais

A Platform nunca executa operações comerciais.

Ela administra o sistema.

---

# Company

Company representa o ambiente de uma empresa.

Cada empresa possui:

- usuários

- pedidos

- estoque

- produtos

- ingredientes

- clientes

- compras

- caixa

- relatórios

Tudo pertence exclusivamente à empresa.

---

# Multi-Tenant

O sistema utiliza isolamento por CompanyID.

Toda entidade pertencente à empresa possui obrigatoriamente:

```
CompanyID
```

Toda consulta deve filtrar por:

```
WHERE company_id = ?
```

Sem exceções.

---

# Usuários da Plataforma

São administradores do SaaS.

Podem:

- criar empresas

- bloquear empresas

- alterar planos

- configurar branding

- visualizar métricas globais

Nunca operam pedidos.

Nunca fazem vendas.

---

# Usuários da Empresa

São funcionários.

Exemplos:

Administrador

Gerente

Caixa

Atendente

Estoquista

Cada um possui permissões específicas.

---

# RBAC

O HorizonGest utiliza RBAC.

Role Based Access Control.

Usuários recebem papéis.

Papéis recebem permissões.

Nunca permissões diretamente ao usuário.

---

# Organização dos módulos

Cada módulo representa um domínio de negócio.

Exemplo

Produtos

Pedidos

Estoque

Compras

Financeiro

Clientes

Relatórios

Cada módulo é independente.

---

# Dependências entre módulos

Permitido

```
Pedido

↓

Estoque
```

Pedido reduz estoque.

---

Também permitido

```
Compra

↓

Estoque
```

Compra aumenta estoque.

---

Não permitido

```
Produto

↓

Financeiro
```

Produto não controla dinheiro.

---

Não permitido

```
Financeiro

↓

Produtos
```

Financeiro não conhece catálogo.

---

# Banco de Dados

Banco único.

Empresas separadas por CompanyID.

Nunca um banco por empresa.

---

# Comunicação

Frontend

↓

HTTP

↓

Handler

↓

Service

↓

Repository

↓

Database

Sempre nesta direção.

---

# Frontend

Responsável apenas por:

- interface

- experiência do usuário

- chamadas HTTP

Nunca possui regra de negócio.

---

# Backend

Responsável por:

- validações

- regras

- segurança

- persistência

- auditoria

Toda inteligência mora aqui.

---

# White Label

O sistema suporta múltiplas marcas.

Branding nunca é fixo.

Nome

Logo

Ícones

Cores

Tudo pode mudar.

---

# Branding

Branding pertence à Platform.

Nunca à empresa.

Uma empresa administra seu negócio.

A Platform administra a identidade do sistema.

---

# Feature Flags

Recursos podem ser ligados ou desligados.

Exemplos

Delivery

Financeiro

Compras

CRM

API Pública

Integrações

Tudo pode ser habilitado sem alterar código.

---

# Segurança

Toda requisição autenticada possui:

UserID

CompanyID

Role

Permissions

Tudo obtido via JWT.

---

# Auditoria

Operações importantes são registradas.

Exemplos

Login

Criação

Alteração

Exclusão

Mudança de plano

Alteração de permissões

Nunca confiar apenas em logs.

---

# Filosofia

O HorizonGest é uma plataforma.

Não é um ERP monolítico.

Não é um sistema para apenas uma empresa.

Tudo foi desenhado para crescer continuamente sem exigir reescrita da arquitetura.

Esse princípio nunca deve ser quebrado.

Fim da Parte 7


# PARTE 8
# Backend do HorizonGest

---

# Capítulo 7
# Arquitetura do Backend

O backend foi construído utilizando Clean Architecture simplificada.

A arquitetura possui separação rígida entre responsabilidades.

Cada camada possui uma única função.

Nenhuma camada conhece detalhes internos das demais.

---

# Estrutura

backend/

├── cmd/

├── internal/

├── migrations/

├── docs/

├── go.mod

└── main.go

---

# cmd/

Contém o ponto de entrada da aplicação.

Exemplo

```
cmd/server/main.go
```

Responsável por:

- iniciar servidor

- carregar configurações

- registrar rotas

- iniciar banco

- iniciar services

- iniciar middlewares

Nunca colocar regra de negócio aqui.

---

# internal/

Toda regra do sistema fica dentro desta pasta.

É dividida em módulos.

---

internal/domain

Modelos de negócio.

Representam conceitos reais do sistema.

Exemplo

Product

Order

Company

Ingredient

Purchase

Theme

BusinessProfile

PlatformBrand

GlobalConfig

Essas estruturas não conhecem banco.

Não conhecem HTTP.

Não conhecem GORM.

São apenas entidades.

---

internal/service

Onde mora TODA regra de negócio.

Esta é a camada mais importante do sistema.

Exemplos

ProductService

OrderService

FinanceService

PurchaseService

AuthService

RBACService

CompanyService

PlatformBrandService

GlobalConfigService

Aqui acontecem:

validações

cálculos

processamentos

regras

transações

integrações

Nunca acessar banco diretamente.

Sempre utilizar Repository.

---

internal/repository

Responsável apenas por persistência.

Nunca faz regra de negócio.

Nunca calcula valores.

Nunca toma decisões.

Sua única função é salvar e buscar dados.

---

internal/handler

Camada HTTP.

Recebe requisição.

Valida entrada.

Chama Service.

Retorna resposta.

Nada além disso.

---

internal/middleware

Executado antes do Handler.

Exemplos

JWT

Tenant

PlatformAuth

RBAC

RateLimit

Logs

Recovery

---

internal/ports

Define interfaces.

Exemplo

```
type ProductRepository interface
```

Service conhece apenas interfaces.

Nunca implementações.

---

internal/infra

Implementações externas.

Banco

SMTP

Storage

Redis

Providers

Tudo que conversa com o mundo externo.

---

# Fluxo de uma requisição

Cliente

↓

HTTP

↓

Router

↓

Middleware

↓

Handler

↓

Service

↓

Repository

↓

Database

Resposta sobe exatamente pelo caminho inverso.

---

# Exemplo

Criar Pedido

↓

POST /orders

↓

OrderHandler

↓

OrderService

↓

OrderRepository

↓

SQLite/Postgres

↓

Response

---

# Handler

O Handler nunca decide nada.

Exemplo correto

```
Receber JSON

↓

Validar campos

↓

Chamar Service

↓

Responder HTTP
```

---

Exemplo incorreto

```
Receber JSON

↓

Calcular desconto

↓

Alterar estoque

↓

Criar pedido

↓

Salvar banco
```

Tudo isso pertence ao Service.

---

# Service

O Service controla o fluxo inteiro.

Exemplo

Criar Pedido

↓

Validar empresa

↓

Validar usuário

↓

Validar estoque

↓

Calcular totais

↓

Gerar número do pedido

↓

Criar pedido

↓

Criar itens

↓

Atualizar estoque

↓

Registrar auditoria

↓

Retornar resultado

Toda inteligência fica aqui.

---

# Repository

Repository apenas executa operações no banco.

Exemplo

```
FindByID()

FindAll()

Save()

Delete()

Update()

Exists()
```

Nunca calcular.

Nunca validar.

Nunca conhecer regra.

---

# Domain

Domain representa o negócio.

Exemplo

Order

não conhece:

SQLite

HTTP

JSON

GORM

Router

Framework

Representa apenas um Pedido.

---

# Dependências

A direção sempre é:

Handler

↓

Service

↓

Repository

↓

Database

Nunca o contrário.

---

# Nunca permitido

Repository chamar Service

Handler chamar Repository

Service chamar Handler

Repository chamar Handler

Domain conhecer banco

Domain conhecer HTTP

Frontend conhecer banco

---

# Interfaces

Todo Service depende de interfaces.

Exemplo

```
ProductRepository
```

e não

```
GormProductRepository
```

Isso permite trocar SQLite por PostgreSQL sem alterar regras.

---

# Banco

Atualmente:

SQLite

No futuro:

PostgreSQL

A arquitetura foi desenhada para que essa troca aconteça alterando apenas a camada Repository.

Nenhuma regra de negócio deve mudar.

---

# Transações

Sempre que várias operações precisam ocorrer juntas, utilizar transação.

Exemplo

Criar Pedido

↓

Criar Itens

↓

Baixar Estoque

↓

Registrar Auditoria

Tudo dentro da mesma transação.

Se qualquer passo falhar:

Rollback.

---

# Concorrência

Operações críticas devem ser preparadas para concorrência.

Exemplo

Geração do número do pedido.

Hoje:

SQLite + Transaction.

Futuro:

PostgreSQL + Advisory Lock.

A regra de negócio permanece exatamente igual.

---

# Filosofia

Cada camada faz apenas uma coisa.

Quanto menor a responsabilidade da camada,

mais fácil será evoluir o sistema.

Nunca adicionar lógica onde ela não pertence.

Fim da Parte 8


# PARTE 9
# Engenharia de Software

---

# Capítulo 8
# Padrões Obrigatórios de Desenvolvimento

Este capítulo define as regras que devem ser seguidas em absolutamente todo novo código desenvolvido para o HorizonGest.

Estas regras possuem prioridade máxima.

Caso exista conflito entre uma decisão local e este documento, este documento prevalece.

---

# Objetivo

Todo código novo deve parecer ter sido escrito pela mesma pessoa.

Isso significa:

- mesma arquitetura

- mesmo padrão

- mesma organização

- mesma nomenclatura

- mesma responsabilidade

A consistência é mais importante que preferência pessoal.

---

# SOLID

Todo desenvolvimento deve respeitar SOLID.

Não por obrigação acadêmica.

Mas porque facilita manutenção.

---

## S

Single Responsibility Principle

Cada classe possui apenas uma responsabilidade.

Exemplo correto

```
OrderService
```

Responsável apenas pelo domínio Pedido.

Não envia e-mail.

Não gera PDF.

Não gera relatórios.

---

Errado

```
OrderService
```

↓

cria pedido

↓

envia email

↓

gera relatório

↓

gera PDF

↓

faz backup

---

## O

Open Closed Principle

O sistema deve crescer adicionando código.

Nunca modificando regras existentes desnecessariamente.

---

## L

Liskov

Implementações devem respeitar interfaces.

Nunca criar exceções escondidas.

---

## I

Interface Segregation

Interfaces pequenas.

Nunca gigantes.

Correto

```
ProductRepository
```

Errado

```
SystemRepository
```

com 300 métodos.

---

## D

Dependency Inversion

Services dependem de interfaces.

Nunca implementações.

---

# DRY

Don't Repeat Yourself.

Nunca duplicar regra de negócio.

Se existe lógica repetida,

extraia.

---

# KISS

Keep It Simple.

Sempre escolher a solução mais simples.

Nunca criar arquitetura complexa sem necessidade.

---

# YAGNI

You Aren't Gonna Need It.

Não implementar funcionalidades para um futuro imaginário.

Só implementar quando existir necessidade real.

---

# Convenção de nomes

Estruturas

Sempre singular.

```
Product

Order

Company

Ingredient
```

Nunca

```
Products

Orders
```

---

Services

Sempre

```
ProductService

OrderService

FinanceService
```

---

Repositories

Sempre

```
ProductRepository

OrderRepository
```

---

Handlers

Sempre

```
ProductHandler

OrderHandler
```

---

Interfaces

Sempre terminam com

Repository

Service

Provider

Storage

Mailer

etc.

---

# Métodos

Utilizar verbos.

Exemplos

```
Create()

Update()

Delete()

Find()

Exists()

Validate()

Calculate()

Generate()

Process()
```

Nunca

```
DoStuff()

Run()

ExecuteThing()
```

---

# Variáveis

Devem explicar exatamente o que representam.

Correto

```
companyID

orderNumber

totalPrice

stockQuantity
```

Errado

```
x

tmp

aux

valor2
```

---

# Comentários

Comentários explicam

POR QUE

Nunca

O QUE

Errado

```go
// soma dois valores
```

Correto

```go
// necessário para manter compatibilidade com pedidos antigos
```

---

# Erros

Sempre encapsular contexto.

Correto

```go
return fmt.Errorf("CreateOrder: gerar número: %w", err)
```

Nunca

```go
return err
```

sozinho.

---

# Logging

Logs precisam ajudar produção.

Correto

```
CompanyID

UserID

OrderID

Operation

Error
```

Nunca

```
Erro desconhecido
```

---

# Panics

Nunca utilizar panic em regra de negócio.

Panic apenas em inicialização irrecuperável.

---

# Transactions

Sempre utilizar transação quando:

uma operação altera múltiplas tabelas.

Exemplos

Pedido

↓

Itens

↓

Estoque

↓

Auditoria

Tudo deve acontecer junto.

---

# Rollback

Qualquer erro cancela toda operação.

Nunca deixar sistema parcialmente atualizado.

---

# Validação

Toda validação pertence ao Service.

Nunca ao Repository.

Nunca ao Frontend.

Frontend apenas melhora UX.

---

# Repository

Repository nunca valida.

Nunca calcula.

Nunca decide.

Repository apenas acessa banco.

---

# Handler

Handler nunca possui inteligência.

Ele apenas:

recebe

↓

valida formato

↓

chama service

↓

retorna resposta

---

# Services

Toda decisão pertence ao Service.

Sempre.

---

# Organização dos arquivos

Um arquivo.

Uma responsabilidade.

Nunca arquivos com milhares de linhas.

Quando crescer demais,

dividir.

---

# Imports

Organizar sempre.

Go padrão.

Bibliotecas.

Projeto.

---

# TODO

Sempre escrever contexto.

Errado

```go
// TODO
```

Correto

```go
// TODO(PostgreSQL): substituir SELECT MAX por Advisory Lock.
```

---

# Refatoração

Sempre preservar comportamento.

Nunca alterar regra de negócio durante refatoração.

Primeiro refatora.

Depois altera comportamento.

---

# Testes

Toda regra importante merece teste.

Especialmente:

autenticação

financeiro

estoque

pedidos

RBAC

branding

multi-tenant

---

# Cobertura

Priorizar testes nas regras críticas.

Não buscar 100% apenas por número.

Buscar qualidade.

---

# Revisão

Antes de aceitar qualquer código perguntar:

Este código quebra arquitetura?

Existe duplicação?

Existe responsabilidade errada?

Existe forma mais simples?

A regra pertence a esta camada?

Se qualquer resposta for "sim",

o código deve ser revisado.

---

# Filosofia

Arquitetura não é decoração.

Ela existe para impedir que o projeto fique impossível de manter daqui a cinco anos.

Toda alteração deve preservar essa filosofia.

Fim da Parte 9


# PARTE 10
# Modelo de Domínio (Domain Model)

---

# Capítulo 9
# Entidades do Sistema

O HorizonGest foi construído utilizando Domain Driven Design simplificado.

Cada entidade representa um conceito real do negócio.

O objetivo do Domain é representar a realidade.

Não representar banco.

Não representar telas.

Não representar API.

---

# Visão Geral

Atualmente o domínio está dividido em grandes grupos.

Platform

↓

Company

↓

Business

↓

Sales

↓

Inventory

↓

Finance

↓

Users

↓

Configuration

---

# PLATFORM

Representa toda a plataforma HorizonGest.

Existe apenas uma plataforma.

Dentro dela existem várias empresas.

A plataforma nunca pertence a uma empresa.

---

## PlatformUser

Representa administradores da plataforma.

Exemplos

Administrador Geral

Suporte

Operações

Financeiro Plataforma

Esses usuários nunca pertencem a Company.

Eles pertencem ao Platform.

---

## PlatformSession

Sessões dos administradores da plataforma.

Usadas para autenticação.

---

## PlatformAudit

Auditoria das operações administrativas.

Nunca misturar com auditoria das empresas.

---

## PlatformBrand

Responsável por toda identidade visual da plataforma.

Exemplos

Nome

Logo

Cor principal

Descrição

Contato

Tema

Toda informação institucional fica aqui.

---

## GlobalConfig

Configurações técnicas globais.

Exemplos

SMTP

Storage

Rate Limit

JWT

Features

Integrações

Nunca armazenar branding aqui.

---

# COMPANY

A Company representa um cliente.

Toda empresa cadastrada no sistema é uma Company.

Exemplos

Padaria

Restaurante

Lanchonete

Pizzaria

Mercado

Cafeteria

Cada Company possui isolamento completo.

---

## Company

Campos importantes

ID

Name

Document

Email

Phone

Status

PlanID

CreatedAt

UpdatedAt

A Company nunca conhece outra Company.

---

## CompanySettings

Configurações específicas da empresa.

Exemplos

Fuso horário

Moeda

Formato fiscal

Impressão

Delivery

PDV

---

## BusinessProfile

Informações institucionais.

Razão Social

Nome Fantasia

Endereço

Cidade

Estado

CEP

CNPJ

Telefone

Logo

Esses dados pertencem somente à empresa.

---

## Theme

Tema visual da empresa.

Cor

Logo

Modo Escuro

Imagem

Personalização

No futuro será utilizado para White Label completo.

---

# USERS

Usuários pertencem sempre a uma Company.

Nunca existem usuários compartilhados entre empresas.

---

## User

Campos principais

ID

CompanyID

Name

Email

PasswordHash

Role

Status

CreatedAt

UpdatedAt

Todo User pertence exatamente a uma Company.

---

## Invitation

Convites enviados para novos usuários.

Fluxo

Empresa

↓

Convida funcionário

↓

Funcionário recebe link

↓

Aceita convite

↓

Conta criada

---

# INVENTORY

Representa estoque.

---

## Category

Categorias de produtos.

Exemplo

Bebidas

Massas

Lanches

Sobremesas

Ingredientes

---

## Ingredient

Matéria-prima.

Exemplos

Farinha

Queijo

Carne

Tomate

Leite

Açúcar

Cada ingrediente possui estoque.

---

## Product

Produto vendido.

Pode utilizar vários ingredientes.

Exemplo

Pizza

↓

Farinha

Molho

Queijo

Calabresa

---

## ProductIngredient

Relacionamento N:N.

Define quais ingredientes fazem parte de um produto.

Também define quantidade utilizada.

---

## StockMovement

Toda alteração no estoque gera movimento.

Entrada

Saída

Cancelamento

Ajuste

Compra

Venda

Nunca alterar estoque sem registrar movimento.

---

## StockAdjustment

Representa ajustes manuais.

Exemplo

Inventário

Perda

Quebra

Correção

---

# PURCHASES

Compras realizadas pela empresa.

---

## Purchase

Compra de fornecedor.

Possui

Fornecedor

Itens

Valor

Status

Data

Empresa

---

## PurchaseItem

Itens da compra.

Cada item aumenta estoque.

---

## PurchaseOrder

Número utilizado pela compra.

Formato atual

PC-{CompanyID}-{Timestamp}

Já é isolado por empresa.

---

# SALES

Representa vendas.

---

## Order

Pedido realizado.

Campos principais

ID

CompanyID

OrderNumber

Status

Customer

Total

Notes

CreatedAt

OrderNumber é sequencial por empresa.

Nunca utilizar ID para exibição.

---

## OrderItem

Itens vendidos.

Cada item referencia Product.

---

## RecentOrder

Modelo utilizado pelo Dashboard.

Serve apenas para consultas.

Não representa tabela física.

---

# FINANCE

Representa fluxo financeiro.

---

## FinanceTransaction

Entrada

Saída

Receita

Despesa

Pagamento

Recebimento

---

## Dashboard

Modelos específicos para consultas.

Nunca utilizados para persistência.

Exemplos

DashboardStats

DashboardSummary

RecentOrders

TopProducts

SalesChart

---

# REPORTS

Modelos específicos para relatórios.

São modelos de leitura.

Nunca entidades persistentes.

---

# MEDIA

Arquivos.

Imagens.

Uploads.

Logos.

Fotos de produtos.

---

# AUTH

Modelos utilizados na autenticação.

JWT

PasswordReset

Blacklist

Sessions

---

# RBAC

Controle de permissões.

Role

Permission

Policies

A autorização sempre acontece no backend.

Nunca confiar no frontend.

---

# RELACIONAMENTOS

Platform

↓

Company

↓

Users

↓

Orders

↓

OrderItems

↓

Products

↓

Ingredients

---

Company

↓

Purchases

↓

PurchaseItems

↓

Ingredients

---

Company

↓

Finance

↓

Dashboard

↓

Reports

---

# REGRAS IMPORTANTES

Toda entidade de negócio possui CompanyID.

Exceto:

PlatformUser

PlatformSession

PlatformAudit

PlatformBrand

GlobalConfig

Plan

Essas pertencem à plataforma.

---

# CompanyID

CompanyID é a principal regra arquitetural do HorizonGest.

Nenhuma consulta pode ignorar CompanyID.

Nenhuma gravação pode omitir CompanyID.

Nenhum Repository pode retornar dados de outra empresa.

---

# Filosofia

O Domain representa o negócio.

Se amanhã SQLite desaparecer...

Se HTTP mudar...

Se React desaparecer...

O Domain continua exatamente igual.

Essa é a principal responsabilidade desta camada.

Fim da Parte 10


# PARTE 11
# Modelo de Dados (Database Model)

---

# Capítulo 10

# Modelo de Banco de Dados

Este documento descreve a estrutura lógica do banco de dados do HorizonGest.

Não descreve SQL.

Não descreve GORM.

Não descreve implementação.

Descreve apenas o modelo de dados.

---

# Filosofia

O banco foi projetado para ser:

• Multi-tenant

• Escalável

• Seguro

• Portável

• Independente do SGBD

Hoje utiliza SQLite.

No futuro poderá utilizar PostgreSQL sem alterar o domínio.

---

# Organização Geral

Platform

↓

Companies

↓

Users

↓

Business

↓

Inventory

↓

Sales

↓

Finance

↓

Reports

---

# CompanyID

A regra mais importante do banco.

Toda tabela de negócio possui:

CompanyID

Sem exceções.

---

Exemplo

Orders

Products

Categories

Ingredients

FinanceTransactions

StockMovements

Purchases

Users

Themes

BusinessProfile

Todos possuem CompanyID.

---

Tabelas globais NÃO possuem CompanyID.

Exemplo

PlatformUser

PlatformAudit

PlatformSession

PlatformBrand

GlobalConfig

Plans

---

# Chaves Primárias

Todas utilizam

ID INTEGER PRIMARY KEY AUTOINCREMENT

Hoje.

No PostgreSQL poderão ser BIGSERIAL.

Ou UUID futuramente.

Essa decisão está isolada na camada Repository.

Nunca no Domain.

---

# IDs

Importante.

O ID interno nunca deve ser mostrado ao usuário.

Ele existe apenas para:

Relacionamentos

Foreign Keys

Referências internas

API

Persistência

---

Sempre que houver necessidade de numeração de negócio deverá existir um campo específico.

Exemplos

OrderNumber

InvoiceNumber

PurchaseNumber

DocumentNumber

Nunca utilizar ID.

---

# Orders

Tabela

orders

Campos principais

ID

CompanyID

OrderNumber

Status

Customer

Notes

TotalPrice

CreatedAt

UpdatedAt

---

OrderNumber

Representa a numeração visível.

É sequencial por empresa.

Empresa A

1

2

3

4

Empresa B

1

2

3

Nunca compartilhar sequência.

---

Índice recomendado

UNIQUE

CompanyID

OrderNumber

---

# OrderItems

Tabela

order_items

Campos

ID

OrderID

ProductID

Quantity

UnitPrice

TotalPrice

CompanyID

Cada item pertence a um pedido.

---

Relacionamentos

Order

↓

OrderItems

↓

Product

---

# Products

Tabela

products

Campos

ID

CompanyID

CategoryID

Name

Description

Price

Active

CreatedAt

---

Relacionamentos

Category

↓

Products

↓

ProductIngredients

---

# Categories

Tabela

categories

Campos

ID

CompanyID

Name

Color

Icon

---

# Ingredients

Tabela

ingredients

Campos

ID

CompanyID

Name

Unit

CurrentStock

MinimumStock

CreatedAt

---

# ProductIngredients

Tabela

product_ingredients

Relacionamento

Product

↓

Ingredient

Possui

Quantidade

Unidade

---

# Purchases

Tabela

purchases

Campos

ID

CompanyID

PurchaseNumber

Supplier

Status

Total

CreatedAt

PurchaseNumber é identificador de negócio.

Nunca utilizar ID.

---

# PurchaseItems

Tabela

purchase_items

Relaciona

Purchase

↓

Ingredient

---

# StockMovements

Tabela

stock_movements

Campos

ID

CompanyID

IngredientID

Type

Quantity

Reason

Reference

CreatedAt

---

Tipos

Purchase

Sale

Adjustment

Loss

Inventory

Transfer

---

Nunca alterar estoque diretamente.

Sempre registrar movimento.

---

# StockAdjustment

Tabela

stock_adjustments

Representa inventários.

Correções.

Perdas.

Acertos.

---

# FinanceTransactions

Tabela

finance_transactions

Campos

ID

CompanyID

Type

Category

Amount

Description

CreatedAt

---

Tipos

Income

Expense

Transfer

Adjustment

---

# Users

Tabela

users

Campos

ID

CompanyID

Name

Email

PasswordHash

Role

Status

---

Email pode repetir entre empresas.

Por isso recomenda-se índice

CompanyID

Email

---

# Invitations

Tabela

invitations

Campos

CompanyID

Email

Token

Expiration

Status

---

# Themes

Tabela

themes

Personalização visual.

CompanyID obrigatório.

---

# BusinessProfile

Tabela

business_profiles

Dados institucionais.

Razão Social

Nome Fantasia

Endereço

Telefone

Logo

CompanyID obrigatório.

---

# PlatformBrand

Tabela

platform_brand

Nunca possui CompanyID.

Existe apenas uma configuração global.

---

# GlobalConfig

Tabela

global_config

Também global.

Nunca possui CompanyID.

---

# PlatformUser

Tabela

platform_users

Administradores.

Nunca pertencem às empresas.

---

# Plans

Tabela

plans

Planos comerciais.

Starter

Professional

Enterprise

White Label

---

# Auditoria

Tabela

platform_audit

Registra operações administrativas.

Não deve registrar operações das empresas.

---

# Sessions

Tabela

platform_sessions

Sessões administrativas.

---

# Password Reset

Tabela

password_resets

Tokens de recuperação.

---

# Token Blacklist

Tabela

token_blacklist

JWT inválidos.

---

# Índices Obrigatórios

Todas as tabelas multi-tenant devem possuir índice em:

CompanyID

---

Índices compostos recomendados

Orders

CompanyID

OrderNumber

---

Users

CompanyID

Email

---

Products

CompanyID

CategoryID

---

Ingredients

CompanyID

Name

---

Finance

CompanyID

CreatedAt

---

Purchases

CompanyID

PurchaseNumber

---

StockMovements

CompanyID

IngredientID

---

# Foreign Keys

Sempre utilizar.

Exemplos

OrderItem

↓

Order

ProductIngredient

↓

Product

ProductIngredient

↓

Ingredient

PurchaseItem

↓

Purchase

StockMovement

↓

Ingredient

---

Nunca utilizar relacionamentos implícitos.

---

# Migrações

Toda alteração estrutural deve ocorrer por Migration.

Nunca alterar banco manualmente.

Cada migration:

é única

imutável

incremental

versionada

---

# Compatibilidade PostgreSQL

O modelo foi desenhado para migrar facilmente.

Diferenças esperadas

SQLite

↓

PostgreSQL

AUTOINCREMENT

↓

BIGSERIAL

TEXT

↓

VARCHAR

BOOLEAN

↓

BOOLEAN

JSON TEXT

↓

JSONB

---

Nenhuma regra de negócio depende do SQLite.

---

# Integridade

Nunca excluir registros importantes.

Preferir:

Status

Active

Soft Delete

Quando aplicável.

---

# Filosofia

O banco existe para armazenar.

Não para decidir regras.

Toda regra permanece no Service Layer.

O banco apenas garante consistência.

Fim da Parte 11.


# PARTE 12
# Arquitetura Backend

---

# Capítulo 11

# Arquitetura Backend

Este documento define a arquitetura oficial do backend do HorizonGest.

Nada pode fugir destas regras.

Qualquer IA ou desenvolvedor deve obedecer este documento.

---

# Objetivos

O backend deve ser:

• simples

• desacoplado

• testável

• escalável

• previsível

• independente do banco

• independente do frontend

---

# Arquitetura Oficial

Sempre:

HTTP

↓

Router

↓

Middleware

↓

Handler

↓

Service

↓

Repository Interface

↓

Repository GORM

↓

Database

---

Nenhum desvio é permitido.

---

# Fluxo completo

Cliente

↓

HTTP Request

↓

Router

↓

Middlewares

↓

Handler

↓

Service

↓

Repository

↓

SQLite/PostgreSQL

↓

Repository

↓

Service

↓

Handler

↓

JSON Response

---

# Camadas

Existem apenas estas camadas.

1

Router

2

Middleware

3

Handler

4

Service

5

Repository

6

Database

---

# Router

Responsável apenas por registrar rotas.

Exemplo

GET

POST

PUT

DELETE

PATCH

Nada mais.

---

Router nunca

• consulta banco

• valida regra

• gera token

• envia email

• calcula estoque

---

# Middleware

Responsável por interceptar requisições.

Exemplos

Autenticação

Autorização

Tenant

Rate Limit

Recovery

Logging

CORS

Headers

---

Middleware nunca

• grava pedido

• altera estoque

• envia email

• salva usuário

---

# Handler

Responsabilidade única.

Traduz HTTP para Service.

Recebe:

Request

↓

Valida entrada

↓

Chama Service

↓

Retorna Response

Nada além disso.

---

Handler pode

✓ ler parâmetros

✓ ler JSON

✓ validar formato

✓ retornar HTTP

---

Handler nunca pode

✗ acessar banco

✗ chamar GORM

✗ fazer SQL

✗ alterar estoque

✗ validar regra de negócio

✗ calcular descontos

✗ gerar números

✗ enviar emails

✗ emitir JWT

---

Exemplo correto

Handler

↓

service.CreateOrder()

↓

return JSON

---

Exemplo errado

Handler

↓

db.Create(...)

---

Errado.

---

Outro exemplo errado

Handler

↓

produto.Stock--

↓

repository.Save()

---

Errado.

---

Toda regra pertence ao Service.

---

# Service

É o cérebro do sistema.

Toda decisão acontece aqui.

Sempre.

---

Service recebe

dados

↓

valida

↓

executa regras

↓

chama repositories

↓

retorna resultado

---

Service pode

✓ validar

✓ calcular

✓ decidir

✓ autorizar

✓ integrar

✓ enviar emails

✓ gerar tokens

✓ chamar APIs

✓ controlar transações

---

Service nunca

✗ recebe HTTP

✗ conhece Gin

✗ conhece Echo

✗ conhece Fiber

✗ escreve JSON

✗ retorna HTTP Status

---

Service nunca sabe que existe navegador.

---

Exemplo

CreateOrder

Service

↓

Validar usuário

↓

Validar empresa

↓

Validar estoque

↓

Gerar OrderNumber

↓

Criar pedido

↓

Baixar estoque

↓

Registrar financeiro

↓

Retornar pedido

Toda lógica fica aqui.

---

# Repository

Responsável apenas por persistência.

Nada além disso.

---

Repository pode

✓ SELECT

✓ INSERT

✓ UPDATE

✓ DELETE

✓ JOIN

✓ paginação

---

Repository nunca

✗ calcula estoque

✗ gera OrderNumber

✗ valida regras

✗ decide autorização

✗ envia email

✗ cria JWT

---

Repository apenas salva.

---

# Repository Interface

Toda comunicação ocorre por interfaces.

Service nunca conhece GORM.

Service conhece apenas interfaces.

Exemplo

OrderRepository

ProductRepository

UserRepository

StockRepository

---

Assim o banco pode mudar.

---

# Repository GORM

Implementação concreta.

Hoje

SQLite

Amanhã

PostgreSQL

Depois

MySQL

Nada muda no Service.

---

# Database

Última camada.

Nunca sobe.

Nunca conhece Service.

Nunca conhece Handler.

---

# Dependências Permitidas

Router

↓

Middleware

↓

Handler

↓

Service

↓

Repository Interface

↓

Repository GORM

↓

Database

A seta sempre aponta para baixo.

---

# Dependências Proibidas

Repository

↓

Service

Nunca.

---

Repository

↓

Handler

Nunca.

---

Service

↓

Handler

Nunca.

---

Database

↓

Service

Nunca.

---

Handler

↓

Repository GORM

Nunca.

---

Handler

↓

Database

Nunca.

---

Service

↓

Database

Nunca.

---

# Comunicação

Sempre unidirecional.

Nunca circular.

---

# Injeção de Dependências

Toda dependência deve ser injetada.

Nunca instanciada dentro do Service.

Correto

main.go

↓

NewOrderService(repo)

---

Errado

service.go

repo := NewRepository()

---

Nunca.

---

# Transações

Toda transação pertence ao Service.

Nunca ao Handler.

---

Exemplo

Criar Pedido

↓

Criar Itens

↓

Atualizar Estoque

↓

Registrar Financeiro

Tudo dentro da mesma transação.

---

# Logs

Logs técnicos

Middleware

Repository

Infrastructure

---

Logs de negócio

Service

---

# Tratamento de Erros

Repository

↓

erro técnico

↓

Service

↓

erro de domínio

↓

Handler

↓

HTTP

Cada camada traduz seus erros.

---

# Retornos

Repository

↓

erro de banco

Service

↓

erro de negócio

Handler

↓

HTTP

---

# Exemplos

Criar Empresa

Handler

↓

Service

↓

Repository

↓

Database

---

Login

Handler

↓

AuthService

↓

UserRepository

↓

JWT

↓

Response

---

Pedido

Handler

↓

OrderService

↓

OrderRepository

↓

StockRepository

↓

FinanceRepository

↓

Response

---

# Benefícios

Arquitetura previsível

Código desacoplado

Testes fáceis

Troca de banco simples

Troca de frontend simples

Baixo acoplamento

Alta manutenção

Escalabilidade

---

# Regra Máxima

Se surgir dúvida:

A decisão sempre pertence ao Service.

Fim da Parte 12.


# PARTE 13
# Arquitetura Frontend

---

# Capítulo 12

# Arquitetura Frontend

O frontend do HorizonGest segue uma arquitetura tão rígida quanto o backend.

Nenhuma IA ou desenvolvedor pode alterar essa estrutura sem justificativa arquitetural.

---

# Objetivos

O frontend deve ser:

• simples

• previsível

• desacoplado

• escalável

• reutilizável

• independente das regras de negócio

---

# Princípio Fundamental

O Frontend NÃO decide regras.

Quem decide é o Backend.

O Frontend apenas apresenta.

---

# Fluxo

Usuário

↓

Interface

↓

Componentes

↓

Store

↓

API

↓

Backend

↓

Resposta

↓

Store

↓

Interface

---

# Camadas

Interface

↓

Componentes

↓

Stores

↓

API Client

↓

Backend

---

# Organização

src/

lib/

components/

stores/

api/

theme/

utils/

routes/

---

# Components

Os componentes possuem responsabilidade única.

Cada componente faz apenas uma coisa.

Exemplos

Button

Modal

Table

Card

Badge

Input

Select

Sidebar

Navbar

Footer

---

Componentes nunca

✗ fazem chamadas HTTP

✗ alteram banco

✗ conhecem regras de negócio

✗ calculam estoque

✗ calculam impostos

---

Componentes apenas exibem dados.

---

# Pages

As páginas apenas organizam componentes.

Nunca concentram lógica.

---

Uma página pode

✓ carregar dados

✓ montar layout

✓ chamar Stores

---

Uma página nunca

✗ calcula regras

✗ conhece banco

✗ possui SQL

---

# Stores

Stores representam estado da aplicação.

Exemplos

AuthStore

BrandStore

ThemeStore

UserStore

CompanyStore

---

Stores podem

✓ armazenar estado

✓ atualizar interface

✓ consumir API

---

Stores nunca

✗ implementam regra de negócio

✗ validam permissões

✗ calculam estoque

---

# API Client

Toda comunicação HTTP passa por um único cliente.

Nunca usar fetch espalhado pelo projeto.

Sempre

Component

↓

Store

↓

API Client

↓

Backend

---

Nunca

Component

↓

fetch()

---

Errado.

---

# Theme

Todo o Design System está centralizado.

theme/

colors

spacing

radius

shadow

typography

animations

transitions

theme.css

---

Nunca utilizar valores mágicos.

Correto

var(--primary)

Errado

#0ea5e9

---

Correto

spacing.md

Errado

margin:17px

---

# Componentes Reutilizáveis

Todo componente deve poder ser reutilizado.

Nunca criar componente exclusivo sem necessidade.

---

Exemplo

Button

↓

usado

Login

Pedidos

Produtos

Financeiro

Configuração

---

# Responsabilidade

Cada componente possui apenas uma responsabilidade.

Exemplo

Button

↓

renderiza botão

Não abre modal.

Não salva pedido.

Não faz login.

---

# Estado

Sempre utilizar Stores.

Nunca espalhar estado entre páginas.

---

# Comunicação

Componentes

↓

Stores

↓

API

Nunca

Componente

↓

Componente

↓

Componente

↓

Componente

---

Evitar cadeia longa de props.

---

# Props

Props apenas para dados.

Nunca enviar funções complexas entre vários níveis.

---

# Eventos

Eventos apenas comunicam intenção.

Exemplo

onClick

↓

Store

↓

Backend

---

# Backend é a Verdade

O Frontend nunca deve assumir sucesso.

Sempre espera resposta.

Backend decide.

Frontend exibe.

---

# Loading

Todo request possui

Loading

Erro

Sucesso

Sempre.

---

# Erros

Todos os erros são tratados.

Nunca deixar Promise sem tratamento.

---

# Organização Visual

Cada módulo possui suas páginas.

Exemplo

Produtos

Pedidos

Financeiro

Estoque

Usuários

Configuração

---

Cada módulo deve permanecer isolado.

---

# Design System

Nunca duplicar CSS.

Sempre utilizar

cores

espaçamentos

bordas

fontes

tokens

---

# Branding

Todo branding vem do backend.

Nunca hardcoded.

Nome

Logo

Cores

Links

Descrição

Tudo dinâmico.

---

# White Label

Todo componente deve funcionar para qualquer cliente.

Nada pode depender do nome HorizonGest.

---

# Responsividade

Toda tela deve funcionar em

Desktop

Tablet

Notebook

Mobile

---

# Acessibilidade

Sempre utilizar

labels

contraste

focus

atalhos

aria

---

# Performance

Evitar renderizações desnecessárias.

Componentes pequenos.

Stores específicas.

Lazy Loading quando necessário.

---

# Organização de Código

Grande

↓

quebrar

Pequeno

↓

manter

Nunca criar páginas enormes.

---

# Regra Máxima

Frontend mostra.

Backend decide.

Sempre.

---

# Benefícios

Código limpo

Baixo acoplamento

Fácil manutenção

Componentes reutilizáveis

Escalabilidade

White Label

Alta produtividade

---

Fim da Parte 13.


# PARTE 14
# Arquitetura Multi-Tenant

---

# Capítulo 13

# Multi-Tenant

O HorizonGest é um sistema SaaS Multi-Tenant.

Toda a arquitetura foi construída considerando que milhares de empresas utilizarão o mesmo sistema simultaneamente.

O isolamento entre empresas é obrigatório.

Nunca opcional.

---

# Definição

Cada empresa enxerga apenas seus próprios dados.

Uma empresa nunca pode acessar informações de outra.

Nem direta.

Nem indiretamente.

Nem através de IDs.

Nem através de relatórios.

Nem através de estatísticas.

---

# Estrutura

Platform

↓

Company A

↓

Dados A

---

Platform

↓

Company B

↓

Dados B

---

Platform

↓

Company C

↓

Dados C

---

Todos compartilham a mesma aplicação.

Todos compartilham o mesmo banco.

Mas nunca compartilham dados.

---

# CompanyID

A principal regra do HorizonGest.

Toda entidade de negócio possui CompanyID.

Sem exceções.

Exemplo

Products

Orders

Ingredients

Purchases

Finance

Users

Themes

BusinessProfile

Categories

StockMovements

Sempre.

---

# Entidades Globais

Apenas poucas entidades NÃO possuem CompanyID.

PlatformUser

PlatformSession

PlatformAudit

PlatformBrand

GlobalConfig

Plans

Essas pertencem à Plataforma.

Nunca às empresas.

---

# Isolamento

Todo acesso ao banco deve possuir CompanyID.

Exemplo

Correto

SELECT *

FROM orders

WHERE company_id = ?

---

Errado

SELECT *

FROM orders

---

Nunca.

---

# Repository

Todo Repository deve filtrar CompanyID.

Mesmo quando o usuário esquecer.

A proteção deve existir na camada Repository.

---

# Services

O Service nunca recebe CompanyID digitado pelo usuário.

Ele recebe CompanyID autenticado.

Nunca confiar em valores enviados pelo frontend.

---

# Middleware

Após autenticação,

o Middleware identifica:

UserID

↓

CompanyID

↓

Role

↓

Permissões

Essas informações acompanham toda a requisição.

---

# JWT

O Token contém:

UserID

CompanyID

Role

Expiration

Issuer

Essas informações nunca devem ser alteradas pelo frontend.

---

# Frontend

O frontend nunca escolhe CompanyID.

Quem define isso é o Backend.

---

# APIs

Toda API autenticada utiliza automaticamente CompanyID.

Nunca permitir:

GET

/orders?company=2

Isso é proibido.

---

O CompanyID vem da autenticação.

Nunca da URL.

---

# Numeração

IDs nunca são exibidos.

Sempre utilizar identificadores de negócio.

Exemplo

OrderNumber

PurchaseNumber

InvoiceNumber

DocumentNumber

---

# Pedido

Empresa A

Pedido 1

Pedido 2

Pedido 3

Empresa B

Pedido 1

Pedido 2

Empresa C

Pedido 1

Cada empresa possui sua própria sequência.

---

# Dashboard

O Dashboard mostra apenas dados da empresa atual.

Nunca dados globais.

---

# Financeiro

Receitas

Despesas

Lucros

Fluxo de Caixa

Tudo filtrado por CompanyID.

---

# Estoque

Cada empresa possui estoque próprio.

Nunca compartilhar ingredientes.

Nunca compartilhar produtos.

---

# Usuários

Usuários pertencem exatamente a uma empresa.

Nunca múltiplas empresas.

Administrador da Plataforma é diferente.

---

# Convites

Convites pertencem à empresa.

Aceitando o convite,

o usuário passa a pertencer apenas àquela Company.

---

# Relatórios

Todos os relatórios devem aplicar filtro CompanyID.

Sem exceção.

---

# Auditoria

Auditoria da Plataforma

↓

PlatformAudit

Auditoria da Empresa

↓

Logs internos da Company

Nunca misturar.

---

# Cache

Sempre considerar CompanyID.

Nunca armazenar cache compartilhado entre empresas.

---

Correto

cache

CompanyID

↓

Produto

---

Errado

cache

↓

Produto

---

# Backup

Backups podem existir em dois níveis.

Empresa

↓

backup da empresa

Plataforma

↓

backup completo

---

# White Label

Cada empresa poderá possuir:

Logo

Cor

Tema

Configuração

Sem afetar as demais.

---

# Segurança

Toda vulnerabilidade Multi-Tenant é considerada crítica.

Nunca aceitar atalhos.

---

# Testes Obrigatórios

Toda funcionalidade nova deve responder:

Empresa A consegue visualizar dados da Empresa B?

Se a resposta puder ser "sim",

há um bug crítico.

---

# Regra de Ouro

Se existir dúvida,

aplique CompanyID.

Se ainda existir dúvida,

aplique novamente.

É melhor filtrar duas vezes do que nenhuma.

---

# Filosofia

O isolamento entre empresas é a característica mais importante do HorizonGest.

Sem Multi-Tenant seguro,

não existe HorizonGest.

Fim da Parte 14.


# PARTE 15
# Autenticação, Autorização e Segurança

---

# Capítulo 14

# Segurança

A segurança do HorizonGest não é opcional.

Ela faz parte da arquitetura.

Nenhum desenvolvedor pode ignorar estas regras.

---

# Objetivos

Garantir

• confidencialidade

• integridade

• disponibilidade

• rastreabilidade

• isolamento entre empresas

---

# Fluxo Geral

Usuário

↓

Login

↓

JWT

↓

Middleware

↓

RBAC

↓

Service

↓

Repository

---

Nenhuma rota protegida pode pular esse fluxo.

---

# Autenticação

O HorizonGest utiliza JWT.

Após login válido:

↓

gera Token

↓

cliente armazena

↓

envia Authorization Bearer

↓

backend valida

↓

permite acesso

---

# JWT

O Token contém apenas:

UserID

CompanyID

Role

Expiration

Issuer

Nunca colocar:

Nome

Email

Senha

Permissões completas

Dados financeiros

---

# Expiração

Todo Token possui validade.

Nunca criar tokens permanentes.

---

# Refresh

Caso exista refresh token,

ele deve possuir tratamento separado.

Nunca reutilizar JWT expirado.

---

# Blacklist

Todo logout invalida o token.

O Token entra na tabela

TokenBlacklist

Mesmo que ainda não tenha expirado.

---

# Password

Senhas nunca são armazenadas.

Apenas PasswordHash.

Nunca texto puro.

Nunca criptografia reversível.

Sempre hash.

---

# Recuperação de Senha

Fluxo

Solicitação

↓

Token único

↓

Validade

↓

Troca senha

↓

Token invalidado

Nunca reutilizar token.

---

# Plataforma

Administradores possuem autenticação independente.

PlatformUser

↓

PlatformSession

Nunca utilizar autenticação das empresas.

---

# Middleware

Toda rota protegida executa:

JWT

↓

CompanyID

↓

Role

↓

RBAC

↓

Handler

---

# RBAC

Role Based Access Control.

Toda autorização acontece no backend.

Nunca confiar no frontend.

---

# Roles

Exemplo

Owner

Manager

Cashier

Kitchen

Employee

PlatformAdmin

Cada role possui permissões específicas.

---

# Permissões

Sempre verificar:

Quem

Pode

Fazer

O quê

Nunca apenas:

Está logado.

---

# Exemplos

Funcionário

↓

não pode excluir empresa

Gerente

↓

não pode acessar plataforma

Administrador Plataforma

↓

não pode acessar dados internos da empresa

---

# CompanyID

Toda autorização considera:

User

↓

Company

↓

Role

↓

Permissão

---

# Isolamento

Mesmo Owner

não pode acessar empresa diferente.

---

# API Pública

Poucas rotas são públicas.

Exemplos

Brand

Health

Login

Forgot Password

Reset Password

Todo restante exige autenticação.

---

# CORS

Sempre configurado.

Nunca utilizar

*

em produção.

---

# Headers

Sempre enviar

Content-Type

Authorization

X-Frame-Options

X-Content-Type-Options

Referrer-Policy

Content-Security-Policy

Quando aplicável.

---

# HTTPS

Produção sempre HTTPS.

Nunca aceitar HTTP.

---

# Rate Limit

Endpoints críticos devem possuir limite.

Exemplo

Login

Reset Password

API Pública

Upload

---

# Upload

Todo upload deve validar

Tipo

Extensão

Tamanho

Nunca confiar apenas no frontend.

---

# SQL Injection

Toda consulta deve utilizar parâmetros.

Nunca concatenar SQL.

---

Correto

WHERE id = ?

---

Errado

WHERE id = " + valor

---

# XSS

Nunca renderizar HTML vindo do usuário.

Sempre escapar conteúdo.

---

# CSRF

Caso utilize cookies,

implementar proteção.

---

# Auditoria

Toda ação administrativa importante deve ser registrada.

Exemplos

Login

Logout

Troca senha

Excluir usuário

Excluir empresa

Alterar plano

Alterar branding

---

# Logs

Nunca registrar

Senha

JWT

Tokens

Cartões

Dados sensíveis

---

# Segredos

Nunca armazenar

JWT Secret

SMTP Password

API Keys

Tokens

no repositório.

Sempre utilizar

Environment Variables.

---

# Variáveis

Exemplos

JWT_SECRET

SMTP_PASSWORD

AWS_SECRET

OPENAI_KEY

Nunca hardcoded.

---

# Criptografia

Sempre utilizar bibliotecas oficiais.

Nunca implementar criptografia própria.

---

# Sessões

Toda sessão deve possuir

Expiração

Revogação

Auditoria

---

# Erros

Nunca revelar detalhes internos.

Correto

Credenciais inválidas.

---

Errado

Senha incorreta.

Usuário inexistente.

Erro SQL.

Stack Trace.

---

# Backup

Backups devem ser protegidos.

Nunca públicos.

Nunca sem criptografia quando armazenados externamente.

---

# White Label

Mesmo utilizando White Label,

a segurança permanece exatamente igual.

---

# Testes

Sempre testar

JWT inválido

JWT expirado

Role incorreta

Empresa diferente

Token revogado

Password Reset expirado

---

# Regra Máxima

Nunca confiar:

no navegador

no frontend

no usuário

na URL

A única fonte de verdade é o Backend.

---

# Filosofia

Segurança não é uma funcionalidade.

É uma característica permanente da arquitetura.

Fim da Parte 15.


# PARTE 16
# Engine de Negócio

---

# Capítulo 15

# Business Engine

O HorizonGest não é um CRUD.

É um motor de negócio.

Toda ação do usuário dispara uma sequência de regras.

Essas regras pertencem exclusivamente ao Service Layer.

Nunca ao Frontend.

Nunca ao Banco.

Nunca ao Handler.

---

# Objetivo

Toda operação deve ser:

determinística

atômica

auditável

consistente

reversível quando possível

---

# Fluxo Geral

Usuário

↓

Handler

↓

Service

↓

Validações

↓

Regras

↓

Repositories

↓

Banco

↓

Resposta

---

Toda regra acontece antes da gravação.

---

# Pedido

Criar um pedido nunca significa apenas gravar uma linha.

Criar um pedido executa um fluxo completo.

---

Fluxo

Receber requisição

↓

Validar usuário

↓

Validar Company

↓

Validar produtos

↓

Validar estoque

↓

Calcular totais

↓

Gerar OrderNumber

↓

Criar Pedido

↓

Criar Itens

↓

Baixar Estoque

↓

Registrar Movimentos

↓

Atualizar Dashboard

↓

Retornar Pedido

---

Nenhuma dessas etapas deve acontecer no Handler.

---

# Estoque

O estoque nunca é alterado diretamente.

Sempre através de um movimento.

---

Fluxo

Venda

↓

StockMovement

↓

Atualiza estoque

↓

Atualiza dashboard

---

Compra

↓

StockMovement

↓

Atualiza estoque

---

Inventário

↓

StockAdjustment

↓

StockMovement

↓

Atualiza estoque

---

Perda

↓

StockAdjustment

↓

StockMovement

↓

Atualiza estoque

---

Nunca executar

ingredient.Stock--

Nunca.

---

# Compras

Compra de fornecedor.

Fluxo

Receber compra

↓

Validar fornecedor

↓

Validar itens

↓

Criar Purchase

↓

Criar PurchaseItems

↓

Gerar PurchaseNumber

↓

Adicionar estoque

↓

Registrar movimentos

↓

Atualizar dashboard

---

# Financeiro

Cada lançamento financeiro possui fluxo próprio.

---

Receita

↓

Validar

↓

Criar transação

↓

Atualizar Dashboard

↓

Retornar

---

Despesa

↓

Validar

↓

Criar transação

↓

Atualizar Dashboard

↓

Retornar

---

# Dashboard

O Dashboard nunca grava dados.

Ele apenas consulta.

Pode utilizar modelos específicos.

Nunca altera banco.

---

# Cadastro de Produto

Fluxo

Validar categoria

↓

Validar empresa

↓

Criar produto

↓

Relacionar ingredientes

↓

Salvar

---

# Cadastro de Ingrediente

Fluxo

Validar unidade

↓

Validar empresa

↓

Criar ingrediente

↓

Criar estoque inicial

↓

Registrar movimento

---

# Usuários

Criar usuário

↓

Validar empresa

↓

Validar e-mail

↓

Gerar senha

↓

Hash

↓

Salvar

↓

Enviar convite

---

# Convites

Empresa

↓

Cria convite

↓

Token

↓

Email

↓

Usuário aceita

↓

Conta criada

↓

Convite encerrado

---

# Login

Receber credenciais

↓

Buscar usuário

↓

Comparar hash

↓

Validar status

↓

Gerar JWT

↓

Registrar sessão

↓

Retornar token

---

# Logout

Receber JWT

↓

Adicionar Blacklist

↓

Encerrar sessão

↓

Retornar sucesso

---

# Recuperação de Senha

Receber email

↓

Gerar token

↓

Salvar token

↓

Enviar email

↓

Usuário redefine

↓

Hash nova senha

↓

Invalidar token

↓

Finalizar

---

# Platform

Fluxos administrativos nunca utilizam Services das empresas.

Sempre possuem Services próprios.

Exemplo

PlatformService

↓

CompanyService

↓

BrandService

↓

PlanService

---

# White Label

Branding

↓

PlatformBrand

↓

API Pública

↓

Frontend

↓

BrandStore

↓

Interface

---

Nenhuma regra de branding pertence ao frontend.

---

# Engine Financeira

No futuro poderá executar automaticamente

Fluxo de Caixa

DRE

CMV

Lucro

Margem

Indicadores

Tudo dentro do Service Layer.

---

# Engine de Estoque

No futuro poderá executar

Reposição

Sugestão de compra

Consumo médio

Validade

Inventário automático

---

# Engine Comercial

No futuro poderá executar

Promoções

Combos

Preço dinâmico

Programa de fidelidade

Cashback

---

# Engine Fiscal

No futuro poderá executar

NF-e

NFC-e

SAT

CF-e

SPED

Sem alterar a arquitetura.

---

# Engine de IA

No futuro poderá executar

Previsão de vendas

Compra automática

Análise financeira

Sugestões de preço

Detecção de perdas

Tudo como novos Services.

Nunca alterando os atuais.

---

# Transações

Sempre que várias operações precisarem ocorrer juntas:

Utilizar transação.

Exemplo

Pedido

↓

Pedido

Itens

Estoque

Financeiro

Dashboard

Tudo confirma.

Ou tudo desfaz.

Nunca parcialmente.

---

# Idempotência

Quando aplicável,

uma mesma operação não deve gerar efeitos duplicados.

Principalmente em integrações externas.

---

# Eventos

Hoje o sistema utiliza chamadas diretas.

No futuro poderá utilizar eventos.

Exemplo

Pedido Criado

↓

Evento

↓

Financeiro

↓

Estoque

↓

Relatórios

↓

Notificações

Sem alterar a lógica principal.

---

# Filosofia

O Business Engine é o coração do HorizonGest.

Toda regra pertence a ele.

Toda evolução futura deve acontecer aqui.

Nunca no Frontend.

Nunca no Banco.

Fim da Parte 16.


# PARTE 17
# Padrões de Desenvolvimento

---

# Capítulo 16

# Development Standards

Este documento define como qualquer código novo deve ser escrito.

Nenhuma IA.

Nenhum desenvolvedor.

Nenhum módulo.

Pode fugir destas regras.

---

# Objetivo

Garantir que o HorizonGest permaneça consistente durante anos.

Mesmo com dezenas de desenvolvedores.

Mesmo utilizando IA.

---

# Filosofia

Sempre priorizar:

clareza

↓

simplicidade

↓

consistência

↓

escalabilidade

↓

performance

---

Nunca sacrificar arquitetura por velocidade.

---

# Estrutura Oficial

Backend

cmd/

internal/

domain/

ports/

service/

handler/

middleware/

infra/

util/

migrations/

---

Frontend

src/

lib/

components/

stores/

api/

theme/

routes/

---

Nunca criar novas estruturas sem necessidade.

---

# Nome dos Arquivos

Sempre utilizar nomes claros.

Exemplo

order_service.go

product_repository.go

company_handler.go

Nunca

service2.go

handler_new.go

teste.go

---

# Nome das Structs

Sempre utilizar nomes completos.

Exemplo

OrderService

ProductRepository

CompanyHandler

Nunca

OS

PR

Srv

---

# Nome dos Métodos

Sempre utilizar verbos.

CreateOrder

UpdateOrder

DeleteOrder

FindOrder

ListOrders

Nunca

Order()

Process()

Execute()

Run()

---

# Services

Cada Service representa um domínio.

Exemplo

OrderService

ProductService

FinanceService

CompanyService

ThemeService

---

Nunca criar

MegaService

SystemService

GeneralService

UtilityService

---

# Repositories

Cada Repository representa apenas uma entidade.

Exemplo

OrderRepository

ProductRepository

IngredientRepository

Nunca

BusinessRepository

DatabaseRepository

MainRepository

---

# Handlers

Cada Handler representa apenas um domínio.

Exemplo

OrderHandler

CompanyHandler

FinanceHandler

Nunca

MainHandler

ApiHandler

SystemHandler

---

# Métodos

Preferir poucos métodos grandes?

Não.

Preferir muitos pequenos.

Cada método deve fazer apenas uma coisa.

---

# Funções

Quando uma função ultrapassar aproximadamente 50 linhas,

avaliar dividir.

Não é regra absoluta.

É orientação.

---

# Comentários

Comentar apenas quando necessário.

Nunca comentar o óbvio.

Correto

// Gera sequência exclusiva por empresa.

Errado

// Soma um ao contador.

contador++

---

# Variáveis

Sempre utilizar nomes claros.

Correto

orderNumber

companyID

ingredientStock

Errado

x

tmp

a

aux

---

# Constantes

Nunca utilizar números mágicos.

Correto

const MaxLoginAttempts = 5

Errado

if attempts > 5

---

# Strings

Evitar repetir textos.

Centralizar quando fizer sentido.

---

# Erros

Sempre encapsular erros.

Exemplo

return fmt.Errorf("CreateOrder: validar estoque: %w", err)

Nunca

return err

---

Assim os logs ficam rastreáveis.

---

# Imports

Sempre remover imports não utilizados.

Nunca deixar código morto.

---

# Código Morto

Nunca comentar código antigo.

Nunca.

Utilizar Git.

---

Errado

// antigo

// db.Create()

novo()

---

Remover.

---

# TODO

Sempre escrever TODO completo.

Exemplo

TODO: Implementar Redis quando houver múltiplas instâncias.

Nunca

TODO

---

# FIXME

Utilizar apenas quando houver bug conhecido.

Nunca utilizar como lembrete.

---

# Logs

Sempre utilizar logs úteis.

Nunca

fmt.Println()

em produção.

---

# Panic

Nunca utilizar panic para regra de negócio.

Panic apenas em falhas irrecuperáveis.

---

# JSON

Sempre utilizar structs.

Nunca map[string]interface{} quando houver modelo conhecido.

---

# SQL

Nunca escrever SQL no Handler.

Nunca.

---

# GORM

Todo GORM permanece na camada Repository.

Nunca subir para Service.

---

# Transactions

Sempre concentradas no Service.

Nunca espalhadas.

---

# Migrations

Uma migration

↓

uma responsabilidade.

Nunca misturar várias alterações grandes.

---

Correto

00031_add_order_number.sql

---

Errado

00031_misc_changes.sql

---

# APIs

Sempre REST.

Recursos claros.

URLs previsíveis.

---

Correto

/orders

/orders/{id}

/products

---

Errado

/doOrder

/process

/createNewOrder

---

# Responses

Sempre retornar objetos consistentes.

Nunca formatos diferentes para a mesma entidade.

---

# Frontend

Nunca duplicar componentes.

Se existir um Button,

utilizar Button.

Nunca criar Button2.

---

# CSS

Nunca utilizar valores soltos.

Sempre Design System.

---

# Cores

Sempre tokens.

Nunca hexadecimal espalhado.

---

# Performance

Primeiro código correto.

Depois otimizar.

Nunca otimizar prematuramente.

---

# Testes

Toda regra nova deve possuir teste.

Principalmente Services.

---

# Cobertura

Prioridade

Service

↓

Repository

↓

Handler

↓

Frontend

---

# Pull Requests

Todo PR deve responder

O que mudou?

Por que mudou?

Impacta arquitetura?

Que testes foram feitos?

---

# IA

Toda IA deve:

ler AI_CONTEXT.md

↓

seguir DECISIONS.md

↓

seguir este documento

↓

produzir código consistente.

---

# Regra Máxima

Se existir duas maneiras de implementar,

escolha sempre a que preserva a arquitetura.

Nunca a mais rápida.

---

# Filosofia

Código bom não é o menor.

Nem o mais inteligente.

Código bom é aquele que outro desenvolvedor consegue entender daqui cinco anos.

Fim da Parte 17.


# PARTE 18
# Testes e Garantia de Qualidade

---

# Capítulo 17

# Testes

O HorizonGest não considera código concluído sem testes.

Testes fazem parte da funcionalidade.

Não são opcionais.

---

# Objetivos

Garantir:

confiabilidade

↓

estabilidade

↓

previsibilidade

↓

segurança

↓

facilidade de manutenção

---

# Pirâmide de Testes

Prioridade

Services

↓

Repositories

↓

Handlers

↓

Frontend

↓

E2E

---

A maior parte dos testes deve existir na camada Service.

É nela que vivem as regras de negócio.

---

# Testes Unitários

Testam apenas uma unidade.

Sem banco.

Sem HTTP.

Sem frontend.

Sempre utilizando mocks quando necessário.

---

Exemplo

OrderService

↓

Mock Repository

↓

Mock StockService

↓

Executa CreateOrder()

↓

Valida resultado

---

# Testes de Repository

Testam acesso ao banco.

Utilizam banco de teste.

Nunca banco de produção.

---

Validam

CRUD

Queries

Filtros

CompanyID

Ordenação

Paginação

---

# Testes de Handler

Testam apenas HTTP.

Request

↓

Middleware

↓

Handler

↓

Response

Nunca testam regra de negócio.

---

# Testes de Frontend

Prioridade menor.

Devem validar:

renderização

componentes

navegação

stores

integração

Nunca lógica de negócio.

---

# Testes E2E

Simulam o usuário.

Fluxo completo.

Exemplo

Login

↓

Criar empresa

↓

Cadastrar ingrediente

↓

Cadastrar produto

↓

Criar pedido

↓

Ver dashboard

---

# Cenários Obrigatórios

Cada funcionalidade nova deve possuir pelo menos:

Fluxo feliz

Erro esperado

Permissão negada

Empresa diferente

Dados inválidos

---

# Testes Multi-Tenant

Sempre obrigatórios.

Pergunta principal:

Empresa A consegue acessar dados da Empresa B?

Se sim,

há falha crítica.

---

Exemplo

Empresa A

Pedido 1

Empresa B

Não pode enxergar Pedido 1.

---

# Testes RBAC

Cada Role deve ser validada.

Owner

Manager

Cashier

Kitchen

Employee

PlatformAdmin

Cada uma deve possuir permissões corretas.

---

# Testes JWT

Validar

Token válido

Token expirado

Token revogado

Token inválido

Issuer inválido

Company diferente

---

# Testes de Estoque

Venda

↓

Baixa estoque

Compra

↓

Aumenta estoque

Ajuste

↓

Atualiza estoque

Perda

↓

Reduz estoque

Nunca permitir estoque inconsistente.

---

# Testes Financeiros

Receita

Despesa

Saldo

Fluxo de Caixa

Devem sempre fechar corretamente.

---

# Testes de Dashboard

Dashboard nunca pode alterar dados.

Somente consultar.

---

# Testes de Concorrência

Sempre que houver geração de sequência.

Exemplo

OrderNumber

PurchaseNumber

Nunca gerar duplicidade.

---

# Testes de Migração

Toda migration deve ser executável.

Banco vazio

↓

Migration

↓

Banco atualizado

↓

Rollback quando aplicável

---

# Testes de Performance

Monitorar

Tempo

Memória

Queries

N+1

Nunca apenas funcionalidade.

---

# Testes de API

Cada endpoint deve validar

200

201

400

401

403

404

409

500

Quando aplicável.

---

# Testes de Upload

Arquivo válido

Arquivo inválido

Extensão inválida

Tamanho inválido

Permissão

---

# Testes de Segurança

JWT

SQL Injection

XSS

CSRF

Rate Limit

Upload

Permissões

Nunca ignorar.

---

# Cobertura

Não perseguimos 100%.

Perseguimos confiança.

Prioridade:

Services

90%+

Repositories

80%+

Handlers

70%+

Frontend

conforme necessidade.

---

# Bug

Todo bug corrigido deve gerar um teste.

Assim ele nunca volta.

---

# Regressão

Quando surgir um bug:

Criar teste

↓

Confirmar falha

↓

Corrigir

↓

Executar teste

Nunca o contrário.

---

# Dados de Teste

Sempre pequenos.

Sempre previsíveis.

Nunca depender de produção.

---

# Seed

O projeto pode possuir seed.

Mas testes nunca dependem do seed.

Cada teste cria seus próprios dados.

---

# Tempo

Testes devem ser rápidos.

Preferencialmente poucos segundos.

Se um teste demora minutos,

algo provavelmente está errado.

---

# CI

Antes de qualquer merge:

Executar

Go Tests

↓

Frontend Check

↓

Build

↓

Lint

↓

Aprovar

---

Nenhum código deve entrar na branch principal quebrado.

---

# Filosofia

Testes não existem para provar que o sistema funciona.

Eles existem para impedir que ele pare de funcionar.

Fim da Parte 18.


# PARTE 19
# Banco de Dados e Migrations

---

# Capítulo 18

# Database Standards

O banco de dados é um dos ativos mais importantes do HorizonGest.

Toda alteração estrutural deve seguir este documento.

Nunca alterar tabelas diretamente em produção.

Nunca alterar estrutura manualmente.

Toda alteração deve passar por migrations.

---

# Filosofia

O banco deve ser:

consistente

↓

versionado

↓

reprodutível

↓

auditável

↓

escalável

---

Toda mudança precisa ser reproduzível em qualquer ambiente.

---

# Banco Oficial

Durante o desenvolvimento:

SQLite

Durante produção:

PostgreSQL

A arquitetura foi construída para funcionar em ambos.

---

# ORM

Utilizamos GORM.

Toda persistência passa pelo Repository Layer.

Nunca utilizar GORM diretamente em Services.

Nunca utilizar GORM em Handlers.

---

# Estrutura

Toda entidade possui:

Domain

↓

Repository

↓

Migration

↓

Service

↓

Handler

---

Nunca criar tabela sem Domain correspondente.

---

# CompanyID

Toda entidade pertencente a uma empresa possui:

CompanyID

Obrigatoriamente.

---

Exemplo

Orders

Products

Ingredients

Categories

Purchases

Finance

Users

Themes

BusinessProfile

---

Nunca esquecer CompanyID.

---

# Entidades Globais

Não possuem CompanyID.

Exemplos

PlatformUser

PlatformAudit

PlatformBrand

GlobalConfig

Plan

PlatformSession

---

Essas entidades pertencem à plataforma.

Nunca às empresas.

---

# IDs

Toda entidade possui

ID

auto incremento

ou UUID futuramente.

Esse ID nunca deve ser mostrado ao usuário.

---

O usuário vê

OrderNumber

PurchaseNumber

InvoiceNumber

etc.

Nunca o ID interno.

---

# Chaves Estrangeiras

Sempre utilizar Foreign Keys.

Sempre que fizer sentido.

Nunca relacionamentos "soltos".

---

# Índices

Criar índices apenas quando necessários.

Prioridade

CompanyID

↓

CreatedAt

↓

Status

↓

Campos de busca

---

Nunca criar índices desnecessários.

---

# Índices Compostos

Sempre avaliar.

Exemplo

(company_id, status)

(company_id, order_number)

(company_id, created_at)

---

Eles costumam melhorar muito consultas.

---

# Nome das Tabelas

Plural.

Exemplo

orders

products

ingredients

companies

users

Nunca

order

produto

tbl_orders

---

# Nome das Colunas

snake_case

Sempre.

---

Exemplo

company_id

order_number

created_at

updated_at

deleted_at

---

Nunca camelCase.

---

# Soft Delete

Sempre utilizar quando fizer sentido.

deleted_at

Permite recuperação.

Preserva histórico.

---

Nunca apagar dados financeiros.

---

# Auditoria

Sempre que necessário,

registrar alterações.

Principalmente

financeiro

estoque

usuários

permissões

configurações

---

# Histórico

Histórico nunca deve depender do log.

Criar entidades próprias quando necessário.

---

# Migration

Cada migration deve possuir apenas uma responsabilidade.

---

Correto

00031_create_order_number.sql

---

Errado

00031_misc.sql

---

# Ordem

As migrations nunca podem ser reorganizadas.

A sequência é definitiva.

---

Nunca alterar migration antiga.

Criar nova.

Sempre.

---

# Alteração de Coluna

Nunca editar uma migration antiga.

Criar

ALTER TABLE

em nova migration.

---

# Remover Colunas

Nunca remover imediatamente.

Fluxo

Adicionar nova

↓

Migrar dados

↓

Atualizar código

↓

Produção

↓

Remover em migration futura

---

# Dados

Nunca colocar dados obrigatórios dentro das migrations.

Exceto

Planos

Permissões

Configurações mínimas

Seeds essenciais

---

# Seed

Seed serve apenas para desenvolvimento.

Nunca depender dele em produção.

---

# Rollback

Sempre que possível,

as migrations devem permitir rollback.

Quando não for possível,

documentar claramente.

---

# Constraints

Sempre utilizar.

NOT NULL

UNIQUE

CHECK

FOREIGN KEY

---

Nunca confiar apenas na aplicação.

---

# UNIQUE

Exemplo

(company_id, order_number)

(company_id, email)

(company_id, sku)

---

Nunca global quando o dado pertence ao tenant.

---

# OrderNumber

OrderNumber nunca utiliza ID global.

Sempre sequência por empresa.

---

SQLite

MAX(order_number)+1

---

PostgreSQL

Preferencialmente

Advisory Lock

ou

Sequence por Tenant

---

Nunca utilizar ID do banco.

---

# Concorrência

Sempre pensar em concorrência.

Principalmente

pedidos

compras

estoque

financeiro

---

Uma operação nunca pode gerar dados duplicados.

---

# Transactions

Toda operação crítica deve utilizar transação.

Exemplo

Pedido

↓

Pedido

Itens

Movimentos

Estoque

Financeiro

Dashboard

Tudo confirma.

Ou tudo desfaz.

---

Nunca deixar gravações parciais.

---

# Queries

Repository é responsável pelas consultas.

Nunca Services escrevem SQL.

Nunca Handlers escrevem SQL.

---

# N+1

Sempre observar.

Se existir:

loop

↓

query

↓

loop

↓

query

Há problema.

Utilizar preload.

Ou JOIN.

---

# Paginação

Toda listagem deve suportar paginação.

Nunca retornar milhares de registros.

---

# Busca

Toda busca deve utilizar índices.

Sempre que possível.

---

# Cache

Banco é fonte oficial.

Cache apenas acelera consultas.

Nunca é verdade absoluta.

---

# Backup

Backups devem ser:

automáticos

↓

versionados

↓

testados

↓

restauráveis

---

Backup não testado não existe.

---

# Restore

Sempre validar restauração.

Não basta gerar backup.

---

# Integridade

Nunca permitir:

pedido sem empresa

produto sem empresa

movimento sem estoque

item sem pedido

---

O banco deve impedir inconsistências.

---

# Performance

Antes de otimizar código,

analisar banco.

Na maioria dos casos,

o gargalo está nas consultas.

---

# Filosofia

O banco de dados deve sobreviver ao código.

Mas o código nunca deve depender de detalhes internos do banco.

Fim da Parte 19.


# PARTE 20
# Deploy, Ambientes e Operação

---

# Capítulo 19

# Operação do Sistema

O HorizonGest foi desenvolvido para funcionar em múltiplos ambientes.

Cada ambiente possui um propósito específico.

Nunca misturar ambientes.

Nunca utilizar produção para testes.

---

# Ambientes Oficiais

Development

↓

Homologação (Staging)

↓

Produção

---

Cada um possui banco independente.

Arquivos independentes.

Configurações independentes.

---

# Development

Objetivo:

Desenvolvimento diário.

Características

SQLite

Logs detalhados

Hot Reload

Dados descartáveis

Sem preocupação com performance máxima.

---

# Homologação

Objetivo

Validar antes da produção.

Características

Configuração semelhante à produção.

Banco próprio.

Testes completos.

Usuários internos.

---

Toda funcionalidade passa por homologação.

---

# Produção

Objetivo

Operação real.

Características

PostgreSQL

Logs controlados

Backup automático

Monitoramento ativo

Alta disponibilidade

---

Produção nunca é laboratório.

---

# Variáveis de Ambiente

Toda configuração sensível deve estar em variáveis.

Nunca no código.

---

Exemplos

DATABASE_URL

JWT_SECRET

SMTP_HOST

SMTP_PASSWORD

SUPABASE_KEY

API_KEYS

---

Nunca versionar segredos.

---

# Arquivos

Exemplo

.env.example

↓

.env.local

↓

.env.production

---

Nunca subir .env real para Git.

---

# Docker

Todo projeto deve subir via Docker.

Componentes

Backend

Frontend

Banco

Redis (quando existir)

---

Docker deve ser suficiente para iniciar todo o ambiente.

---

# Docker Compose

Responsável por integrar todos os serviços.

Nunca depender de instalação manual.

---

# Build

Backend

↓

Go Build

↓

Executável

---

Frontend

↓

npm run build

↓

Artefato final

---

Nenhum build pode depender de arquivos locais.

---

# Deploy

Fluxo ideal

Merge

↓

CI

↓

Build

↓

Testes

↓

Deploy

↓

Health Check

---

Nunca fazer deploy manual copiando arquivos.

---

# CI/CD

Todo merge para main deve executar

Lint

↓

Go Test

↓

Svelte Check

↓

Build Backend

↓

Build Frontend

↓

Deploy (quando configurado)

---

Se algum passo falhar,

o deploy não acontece.

---

# Logs

Logs devem ser úteis.

Nunca excessivos.

---

Devem conter

Data

Hora

Nível

Serviço

Mensagem

Contexto

---

Nunca registrar senhas.

Nunca registrar tokens.

Nunca registrar dados sensíveis.

---

# Níveis

DEBUG

INFO

WARN

ERROR

FATAL

---

Produção normalmente utiliza

INFO

WARN

ERROR

---

# Monitoramento

Monitorar

CPU

Memória

Banco

Tempo de resposta

Erros

Fila

Backup

---

Nunca esperar o cliente descobrir problemas.

---

# Health Check

Todo serviço deve possuir

/health

---

Retorna

Status

Banco

Versão

Tempo

Dependências

---

Exemplo

200 OK

{
status:"healthy"
}

---

# Graceful Shutdown

Ao desligar o servidor

Finalizar requisições

↓

Fechar conexões

↓

Salvar estados necessários

↓

Encerrar

Nunca interromper operações críticas.

---

# Backup

Automático.

Versionado.

Testado.

---

Backup mínimo

Diário

Retenção

30 dias

---

Produção crítica

Hora em hora.

---

# Restore

Backup deve ser restaurável.

Sempre validar.

Backup sem restore testado não possui valor.

---

# Atualizações

Atualizações seguem

Migration

↓

Deploy Backend

↓

Deploy Frontend

↓

Validação

↓

Liberação

Nunca alterar banco depois do deploy.

---

# Rollback

Toda atualização deve possuir estratégia de retorno.

Código

Banco

Arquivos

Configuração

---

Se necessário,

rollback deve ocorrer rapidamente.

---

# Banco

Nunca acessar banco de produção manualmente.

Exceto situações extremamente justificadas.

---

Toda alteração estrutural passa por migration.

---

# Cache

Pode ser limpo.

Pode ser recriado.

Nunca armazenar dados únicos apenas no cache.

---

# Arquivos

Uploads devem possuir backup próprio.

Banco e arquivos possuem ciclos independentes.

---

# Segurança

Sempre utilizar

HTTPS

↓

JWT

↓

Headers

↓

Rate Limit

↓

Validação

---

Nunca expor endpoints administrativos publicamente.

---

# Observabilidade

O sistema deve responder rapidamente

O que aconteceu?

Quando aconteceu?

Quem executou?

Onde ocorreu?

---

Isso reduz drasticamente tempo de suporte.

---

# Escalabilidade

Toda arquitetura deve permitir

Mais usuários

↓

Mais empresas

↓

Mais módulos

↓

Mais servidores

Sem reescrever o sistema.

---

# Deploy Zero Downtime

Objetivo futuro.

Atualização sem interromper usuários.

Arquitetura preparada para isso.

---

# Filosofia

Deploy não é o fim do desenvolvimento.

É o início da operação.

A qualidade de um sistema é medida muito mais pela estabilidade em produção do que pela velocidade com que foi escrito.

Fim da Parte 20.


# PARTE 21
# Evolução do Produto e Governança

---

# Capítulo 20

# Product Governance

O HorizonGest foi construído para evoluir continuamente.

Sua arquitetura não foi criada apenas para atender às necessidades atuais.

Ela foi projetada para suportar anos de crescimento.

Toda evolução do produto deve respeitar este documento.

---

# Objetivo

Garantir que:

novas funcionalidades

↓

novos módulos

↓

novas integrações

↓

novas tecnologias

possam ser adicionadas

sem quebrar a arquitetura existente.

---

# Regra Principal

A arquitetura é permanente.

As funcionalidades evoluem.

Nunca o contrário.

---

# Toda Nova Funcionalidade

Antes de ser implementada deve responder:

Qual problema resolve?

Qual módulo pertence?

Precisa de novo domínio?

Precisa de nova entidade?

Precisa alterar arquitetura?

Se alterar arquitetura,

deve ser discutido antes.

---

# Regra de Ouro

Nunca modificar código existente apenas para encaixar uma nova funcionalidade.

Preferir estender.

Nunca deformar.

---

# Novo Módulo

Todo módulo novo deve possuir:

Domain

↓

Repository

↓

Service

↓

Handler

↓

Frontend

↓

Documentação

↓

Testes

---

Nunca criar módulos "pela metade".

---

# Exemplo

CRM

Financeiro

Produção

Delivery

Compras

Estoque

RH

Todos seguem exatamente o mesmo padrão.

---

# Novas Integrações

Toda integração externa deve possuir um Adapter próprio.

Nunca misturar regras da integração dentro do Service principal.

---

Exemplo

WhatsApp

PIX

OpenAI

Mercado Pago

iFood

Cada um possui sua própria camada.

---

# IA

Qualquer inteligência artificial futura

não deve alterar regras existentes.

Ela apenas auxilia.

Exemplo

Sugestões

Previsões

Análises

Alertas

Nunca decisões automáticas sem autorização explícita.

---

# Feature Flags

Toda funcionalidade grande deve poder ser ligada ou desligada.

Isso facilita:

implantação

testes

rollback

clientes específicos

---

# White Label

Toda funcionalidade nova deve respeitar o White Label.

Nunca criar:

logos fixos

cores fixas

nomes fixos

URLs fixas

---

Tudo deve consumir:

PlatformBrand

BrandStore

Configurações da Plataforma

---

# Compatibilidade

Sempre que possível,

novas versões devem manter compatibilidade com versões anteriores.

Principalmente na API.

---

# API

Nunca quebrar endpoints existentes.

Se necessário,

criar nova versão.

Exemplo

/api/v1

/api/v2

---

# Banco

Mudanças estruturais seguem migrations.

Nunca alterar banco manualmente.

---

# Performance

Toda funcionalidade nova deve responder:

Quanto consome?

Quanto consulta?

Quanto grava?

Quanto escala?

---

# Escalabilidade

Antes de implementar,

perguntar:

Funciona com

10 empresas?

100 empresas?

1.000 empresas?

10.000 empresas?

---

Se a resposta for não,

repensar.

---

# Observabilidade

Toda funcionalidade crítica deve gerar logs suficientes para investigação futura.

---

# Auditoria

Toda alteração administrativa importante deve ser auditada.

Nunca perder histórico.

---

# Roadmap

O Roadmap técnico possui prioridade sobre desejos momentâneos.

Nunca criar funcionalidades improvisadas apenas porque "é rápido".

---

# Refatoração

Refatoração é permitida.

Mas deve preservar comportamento.

Nunca alterar regra de negócio durante refatoração.

---

# Dívida Técnica

Pode existir.

Mas deve ser registrada.

Nunca escondida.

Todo TODO importante deve possuir contexto.

---

# Revisão

Toda alteração estrutural deve responder:

Melhora a arquitetura?

Ou apenas resolve um problema imediato?

Se for apenas imediatismo,

normalmente não deve ser aceita.

---

# Inteligência Artificial

Toda IA que trabalhar no projeto deve:

Ler AI_CONTEXT.md

↓

Ler DECISIONS.md

↓

Ler este Manual Mestre

↓

Seguir a arquitetura

Nunca criar sua própria arquitetura.

---

# Filosofia

O HorizonGest não deve crescer por acúmulo de código.

Ele deve crescer por evolução organizada.

Cada novo módulo deve parecer que sempre fez parte do sistema.

Essa é a diferença entre um software que envelhece bem e um software que precisa ser reescrito.

Fim da Parte 21.


# PARTE 22
# Como Criar um Novo Módulo

---

# Capítulo 21

# Guia Oficial de Criação de Módulos

Este capítulo define o único processo aceito para adicionar funcionalidades ao HorizonGest.

Todo novo módulo deve seguir exatamente este fluxo.

Não existem exceções.

---

# Objetivo

Garantir que todos os módulos:

sigam a arquitetura

↓

tenham baixa manutenção

↓

possam evoluir

↓

não gerem acoplamento

↓

sejam previsíveis

---

# O Fluxo

Antes de escrever código:

Entender o problema

↓

Modelar o domínio

↓

Criar entidades

↓

Criar Repository

↓

Criar Service

↓

Criar Handler

↓

Criar Frontend

↓

Criar testes

↓

Documentar

Somente depois considerar o módulo concluído.

---

# Passo 1

Entender o Domínio

Nunca começar escrevendo código.

Primeiro entender:

Quem usa?

Qual problema resolve?

Quais regras existem?

Quais entidades serão criadas?

Quais relacionamentos existem?

---

# Passo 2

Criar o Domain

Toda funcionalidade começa aqui.

Criar:

Struct

Validações básicas

Interfaces

Tipos

Constantes

Nunca colocar acesso ao banco.

Nunca colocar HTTP.

Nunca colocar frontend.

---

# Passo 3

Criar Repository

Responsável exclusivamente por persistência.

Implementar:

Create

Update

Delete

Find

List

Search

Paginação

Filtros

CompanyID

Nada além disso.

---

# Passo 4

Criar Service

Aqui vive toda regra de negócio.

Exemplos

validações

regras

cálculos

estoque

financeiro

permissões

integrações

Nunca acessar banco diretamente.

Sempre utilizar Repository.

---

# Passo 5

Criar Handler

Responsável apenas por HTTP.

Recebe Request.

↓

Valida Input.

↓

Chama Service.

↓

Retorna Response.

Nunca escrever regra de negócio.

---

# Passo 6

Registrar Rotas

Adicionar apenas as rotas necessárias.

Sempre protegidas pelos middlewares corretos.

Nunca criar endpoint público sem necessidade.

---

# Passo 7

Criar Frontend

O frontend apenas consome API.

Fluxo

Página

↓

Store (quando necessário)

↓

API Client

↓

Backend

Nunca implementar regra de negócio.

---

# Passo 8

Atualizar Navegação

Se o módulo possuir interface:

Sidebar

Menus

Breadcrumb

Permissões

Sempre respeitando RBAC.

---

# Passo 9

Criar Testes

Obrigatório.

Service

Repository

Handler

Fluxo principal

Empresa diferente

Permissões

Erros

Sem testes o módulo não está finalizado.

---

# Passo 10

Atualizar Documentação

Todo módulo deve atualizar:

AI_CONTEXT.md

↓

README do módulo

↓

API

↓

Roadmap (quando necessário)

---

# Estrutura Esperada

Exemplo

internal/domain/customer.go

internal/ports/customer_repository.go

internal/infra/repository/gorm_customer_repository.go

internal/service/customer_service.go

internal/handler/customer_handler.go

frontend/routes/customers/

frontend/lib/api/customer.ts

docs/

tests/

---

Sempre esse padrão.

---

# CompanyID

Se o módulo pertence à empresa:

CompanyID obrigatório.

Nunca esquecer.

---

# Auditoria

Pergunta obrigatória:

Este módulo precisa gerar histórico?

Se sim,

implementar auditoria.

---

# Permissões

Pergunta obrigatória:

Quem pode acessar?

Owner?

Manager?

Cashier?

Kitchen?

Employee?

Platform?

Nunca deixar implícito.

---

# Dashboard

Pergunta obrigatória.

Este módulo impacta indicadores?

Se sim,

atualizar Dashboard.

---

# Financeiro

Pergunta obrigatória.

Este módulo movimenta dinheiro?

Se sim,

registrar movimentação financeira.

Nunca atualizar saldo diretamente.

---

# Estoque

Pergunta obrigatória.

Este módulo altera estoque?

Se sim,

registrar movimento.

Nunca alterar quantidade diretamente.

---

# Eventos

Sempre pensar:

O que aconteceu?

Pedido criado?

Compra cancelada?

Produto alterado?

Esses eventos poderão ser úteis futuramente.

---

# Integrações

Se o módulo conversar com sistemas externos:

Criar Adapter.

Nunca misturar lógica da integração.

---

# Cache

Perguntar:

Precisa?

Se não,

não implementar.

Cache é otimização.

Nunca requisito.

---

# API

Todo endpoint deve possuir:

Request

Response

Erros

Permissões

Documentação

---

# UI

Toda tela nova deve seguir:

Design System

↓

Theme

↓

Components

↓

Spacing

↓

Typography

Nunca criar componentes "fora do padrão".

---

# Checklist Final

Antes de concluir:

☐ Domain criado

☐ Repository criado

☐ Service criado

☐ Handler criado

☐ Rotas registradas

☐ Frontend criado

☐ RBAC atualizado

☐ Dashboard atualizado (se necessário)

☐ Financeiro atualizado (se necessário)

☐ Estoque atualizado (se necessário)

☐ Testes criados

☐ Documentação atualizada

Somente após todos os itens marcados o módulo pode ser considerado concluído.

---

# Filosofia

Um módulo novo deve parecer que sempre existiu.

Se ao olhar o código alguém consegue identificar exatamente quando ele foi criado,

provavelmente ele não respeitou a arquitetura.

Fim da Parte 22.


# PARTE 23
# Manual Oficial do Arquiteto de Software

---

# Capítulo 22

# O Papel do Arquiteto

O Arquiteto não é o desenvolvedor mais experiente.

O Arquiteto é o guardião da arquitetura.

Sua responsabilidade principal é garantir que o HorizonGest continue evoluindo sem perder consistência.

---

# Missão

Preservar:

Arquitetura

↓

Qualidade

↓

Escalabilidade

↓

Padronização

↓

Visão de longo prazo

---

# O que o Arquiteto NÃO faz

Não programa tudo.

Não centraliza conhecimento.

Não cria dependências pessoais.

Não toma decisões por preferência.

Não escolhe tecnologias apenas porque são novas.

---

# O que o Arquiteto faz

Define padrões.

Revisa arquitetura.

Avalia impacto.

Controla dívida técnica.

Protege a fundação.

Orienta a equipe.

---

# A Regra Principal

Toda decisão deve favorecer o projeto.

Nunca favorecer o desenvolvedor.

Nunca favorecer uma tecnologia.

Nunca favorecer velocidade momentânea.

---

# Antes de Aprovar uma Mudança

Responder:

Resolve um problema real?

É consistente com a arquitetura?

Escala?

É simples?

Pode ser mantida daqui cinco anos?

---

Se qualquer resposta for "não",

a mudança deve ser revista.

---

# Arquitetura

A arquitetura é mais importante que qualquer módulo.

Se um módulo exigir quebrar a arquitetura,

o módulo deve mudar.

Nunca a arquitetura.

---

# Simplicidade

Sempre escolher a solução mais simples

que resolva completamente o problema.

Complexidade sem necessidade é defeito.

---

# Consistência

É melhor dez soluções iguais

do que dez soluções "mais inteligentes".

Consistência reduz manutenção.

---

# Tecnologia

Tecnologia não é objetivo.

É ferramenta.

Nunca trocar tecnologia apenas porque surgiu algo novo.

---

# Dependências

Antes de adicionar qualquer biblioteca perguntar:

Realmente precisamos?

Já existe solução interna?

Quanto isso aumenta a complexidade?

Quem manterá isso?

---

Cada dependência deve justificar sua existência.

---

# Refatoração

Refatorar é permitido.

Mas somente quando:

melhora a arquitetura

↓

reduz complexidade

↓

reduz duplicação

↓

aumenta clareza

Nunca refatorar apenas por gosto pessoal.

---

# Código

Todo código deve ser:

legível

↓

previsível

↓

testável

↓

documentado

↓

simples

---

# Escalabilidade

Sempre pensar além do cenário atual.

Hoje:

2 empresas.

Amanhã:

20.

Depois:

2.000.

Arquitetura não deve mudar.

---

# Multi-Tenant

É inegociável.

Toda decisão deve preservar:

isolamento

↓

segurança

↓

CompanyID

↓

independência

---

# White Label

É parte da fundação.

Nunca criar:

logos fixos

cores fixas

nomes fixos

textos institucionais fixos

---

# Banco

Nunca permitir atalhos.

Toda alteração:

migration

↓

repository

↓

service

↓

handler

Nunca SQL espalhado.

---

# Frontend

Nunca colocar regra de negócio.

Nunca acessar banco.

Nunca criar estados duplicados.

---

# Backend

Service contém regras.

Repository contém persistência.

Handler contém HTTP.

Nunca misturar responsabilidades.

---

# Performance

Nunca otimizar antes de medir.

Primeiro medir.

Depois otimizar.

---

# Segurança

Segurança nunca é opcional.

Toda funcionalidade deve responder:

Existe risco?

Pode expor dados?

Pode quebrar isolamento?

Pode aumentar superfície de ataque?

---

# Documentação

Arquitetura sem documentação não existe.

Toda decisão importante deve ser registrada.

---

# Revisão de Código

Toda revisão deve responder:

Está correto?

Está simples?

Está consistente?

Está seguro?

Está alinhado ao manual?

---

# Dívida Técnica

Pode existir.

Mas nunca escondida.

Toda dívida técnica deve estar documentada.

---

# IA

A IA auxilia.

O Arquiteto decide.

A responsabilidade final sempre é humana.

---

# Evolução

O objetivo do HorizonGest não é crescer rapidamente.

É crescer corretamente.

Velocidade sem direção produz caos.

Arquitetura sem evolução produz obsolescência.

O equilíbrio é responsabilidade do Arquiteto.

---

# Filosofia

O Arquiteto deve tomar decisões pensando na próxima década,

não na próxima sprint.

Uma boa arquitetura permite que centenas de funcionalidades sejam adicionadas sem perder organização.

Essa é a verdadeira medida do sucesso arquitetural.

Fim da Parte 23.


# PARTE 24
# Manual Oficial das Inteligências Artificiais do HorizonGest

---

# Objetivo

Este documento define como qualquer Inteligência Artificial deve trabalhar dentro do HorizonGest.

Inclui:

ChatGPT

Cascade

Claude

Gemini

Copilot

Cursor

Windsurf

Codex

qualquer IA futura.

Nenhuma IA pode ignorar estas regras.

---

# A IA nunca é dona do projeto

A IA é uma ferramenta.

Nunca decide arquitetura.

Nunca decide engenharia.

Nunca muda padrões.

Nunca muda convenções.

Ela apenas executa aquilo que o Manual Oficial determina.

---

# A IA nunca cria arquitetura nova

É proibido:

inventar camadas

inventar padrões

inventar organização

inventar estruturas

inventar modelos

A arquitetura oficial já existe.

---

# A IA nunca pode simplificar regras

Mesmo que ache melhor.

Mesmo que conheça outro padrão.

Mesmo que considere mais moderno.

O HorizonGest possui arquitetura própria.

Ela deve ser seguida.

---

# Antes de escrever código

A IA deve responder internamente:

Estou respeitando a arquitetura?

Estou respeitando os módulos?

Estou respeitando o isolamento?

Estou respeitando os Services?

Estou respeitando os Repositories?

Estou respeitando os Handlers?

Se alguma resposta for negativa,

o código não deve ser produzido.

---

# A IA nunca muda padrões existentes

Se existe um padrão para:

Repository

Service

Handler

Migration

Model

DTO

Controller

Middleware

Store

Componente

ela deve reutilizar exatamente o mesmo padrão.

---

# Proibido criar "atalhos"

Nunca:

acessar banco diretamente

consultar SQL dentro do Service

misturar Repository com Service

misturar Frontend com Backend

misturar regra de negócio com interface

---

# Toda regra pertence ao Service

Sempre.

Sem exceções.

---

# Toda persistência pertence ao Repository

Sempre.

Sem exceções.

---

# Todo HTTP pertence ao Handler

Sempre.

Sem exceções.

---

# Frontend

O frontend nunca contém regra de negócio.

Pode:

mostrar

validar UX

formatar

animar

renderizar

Nunca:

calcular negócio

tomar decisões de negócio

consultar banco

duplicar regras do backend

---

# Banco de Dados

Toda alteração exige Migration.

Nunca alterar tabelas manualmente.

Nunca alterar estrutura diretamente.

Nunca assumir estado do banco.

---

# Multi-Tenant

É obrigatório preservar:

CompanyID

isolamento

segurança

independência

Toda consulta deve considerar isso.

---

# IDs

A IA deve distinguir:

ID Global

↓

uso interno

e

Número de Negócio

↓

visível ao usuário

Jamais mostrar IDs internos quando existir identificador funcional.

---

# White Label

É proibido escrever:

"HorizonGest"

"PratoOnline"

logos

cores

nomes

textos institucionais

hardcoded.

Tudo deve vir do Branding.

---

# Branding

Sempre consumir:

PlatformBrand

BrandStore

BrandService

Nunca escrever valores fixos.

---

# Código duplicado

Antes de criar código,

a IA deve procurar se já existe.

Se existir,

reutilizar.

Nunca duplicar.

---

# Refatoração

A IA nunca refatora por iniciativa própria.

Somente quando solicitado.

---

# Dependências

Nunca instalar bibliotecas sem autorização explícita.

Primeiro justificar.

Depois solicitar aprovação.

---

# Arquivos

Nunca criar arquivos desnecessários.

Antes perguntar:

Existe arquivo para isso?

Existe módulo para isso?

Existe componente equivalente?

---

# Nomeação

Seguir exatamente o padrão existente.

Nunca misturar:

snake_case

camelCase

PascalCase

kebab-case

Cada camada possui seu padrão.

---

# Comentários

Só comentar quando necessário.

Código deve explicar a si mesmo.

---

# Logs

Logs devem ser úteis.

Nunca verbosos.

Nunca expor dados sensíveis.

Nunca expor tokens.

Nunca expor senhas.

---

# Segurança

A IA deve assumir que qualquer entrada é maliciosa.

Sempre validar.

Sempre sanitizar.

Sempre proteger.

---

# Performance

Nunca otimizar prematuramente.

Primeiro resolver corretamente.

Depois otimizar.

---

# Testes

Toda funcionalidade importante deve possuir testes.

Nunca remover testes existentes.

Nunca quebrar testes existentes.

---

# Documentação

Toda decisão arquitetural importante deve atualizar:

DECISIONS.md

quando solicitado.

Toda mudança estrutural deve atualizar:

AI_CONTEXT.md

quando necessário.

---

# Git

Nunca criar commits automaticamente.

Nunca decidir mensagens de commit.

Sempre deixar o commit para o desenvolvedor.

---

# Mudanças grandes

Mudanças estruturais devem ser divididas em pequenas etapas.

Nunca alterar dezenas de arquivos sem justificativa.

---

# Comunicação

A IA deve responder de forma objetiva.

Quando existir risco,

explicar.

Quando existir alternativa,

mostrar.

Quando existir impacto,

informar.

---

# Em caso de dúvida

A IA deve perguntar.

Nunca assumir.

Nunca inventar.

---

# Prioridades

Sempre obedecer esta ordem:

1 Arquitetura

↓

2 Segurança

↓

3 Isolamento Multi-Tenant

↓

4 Integridade dos Dados

↓

5 Clareza

↓

6 Performance

↓

7 Conveniência

---

# Filosofia

A IA existe para acelerar o desenvolvimento,

não para alterar a identidade técnica do HorizonGest.

Toda sugestão deve fortalecer a arquitetura existente.

Nunca substituí-la.

---

# Regra Final

Se qualquer resposta produzida pela IA violar qualquer documento oficial do projeto,

a resposta deve ser considerada incorreta,

mesmo que funcione tecnicamente.

Arquitetura sempre vence.

Fim da Parte 24.


# PARTE 25

# Checklist Obrigatório antes de qualquer Merge

Todo código desenvolvido para o HorizonGest deve passar por esta lista antes de ser considerado pronto.

---

## Arquitetura

☐ Seguiu o ARCHITECTURE_RULES.md

☐ Não criou novas camadas

☐ Não alterou arquitetura

☐ Não criou dependências circulares

☐ Não quebrou desacoplamento

---

## Backend

☐ Handler apenas HTTP

☐ Service apenas negócio

☐ Repository apenas persistência

☐ Sem SQL fora do Repository

☐ Sem lógica de negócio no Handler

---

## Frontend

☐ Sem regra de negócio

☐ Apenas UX

☐ Apenas renderização

☐ Sem duplicar validações do backend

---

## Banco

☐ Toda alteração possui Migration

☐ Nenhuma alteração manual

☐ Nenhuma tabela criada fora do padrão

---

## Multi-Tenant

☐ CompanyID respeitado

☐ Queries isoladas

☐ Nenhum vazamento entre empresas

☐ Nenhuma consulta global indevida

---

## Segurança

☐ Inputs validados

☐ Outputs seguros

☐ Tokens protegidos

☐ Logs sem informações sensíveis

☐ JWT preservado

---

## White Label

☐ Sem HorizonGest hardcoded

☐ Sem PratoOnline hardcoded

☐ Sem logos fixos

☐ Sem cores fixas

☐ Branding consumido via serviço

---

## Performance

☐ Sem consultas N+1

☐ Sem processamento duplicado

☐ Sem carregamentos desnecessários

☐ Sem loops evitáveis

---

## Código

☐ Sem duplicação

☐ Sem código morto

☐ Sem TODO esquecido

☐ Sem comentários antigos

☐ Sem debug

☐ Sem console.log

☐ Sem fmt.Println

---

## Testes

☐ Testes existentes continuam passando

☐ Novas regras possuem testes

☐ Nenhum teste removido

---

## Documentação

☐ AI_CONTEXT atualizado (quando necessário)

☐ DECISIONS atualizado (quando necessário)

☐ Manual atualizado (quando necessário)

---

# Critério de Aceite

Uma funcionalidade somente pode ser considerada pronta quando todos os itens acima forem verdadeiros.

Caso contrário,

ela permanece em desenvolvimento.

---

# Definição Oficial de Pronto

Uma funcionalidade está pronta quando:

✔ Funciona

✔ Está documentada

✔ Possui testes

✔ Respeita arquitetura

✔ Respeita engenharia

✔ Respeita segurança

✔ Respeita isolamento

✔ Não cria dívida técnica

---

# Regra Final

No HorizonGest:

"Funcionar" nunca é suficiente.

"Funcionar corretamente dentro da arquitetura" é o único critério válido.

Fim da Parte 25.
```


# PARTE 26

# Roadmap Oficial do Produto

Este documento representa a ordem oficial de evolução do HorizonGest.

Ela deve ser seguida por qualquer desenvolvedor ou IA.

Mudanças de ordem somente podem ocorrer mediante decisão arquitetural registrada no DECISIONS.md.

---

# Fase 1 — Fundação (CONCLUÍDA)

Status:
✅ Concluída

Objetivos:

- Arquitetura
- Engenharia
- Multi-tenant
- Segurança
- White Label
- RBAC
- Plataforma
- Documentação
- Manual
- Testes

Resultado:

Foundation Closed

Nota:

9.0/10

---

# Fase 2 — Estabilização

Status:
Em andamento

Objetivo:

Executar o sistema do zero até operação completa.

Corrigir:

- bugs
- fluxo
- UX
- integrações
- permissões
- inconsistências

Nenhuma funcionalidade nova entra nesta fase.

Somente correções.

Checklist:

□ Cadastro Plataforma

□ Cadastro Empresa

□ Convites

□ Funcionários

□ Ingredientes

□ Produtos

□ Compras

□ Estoque

□ Ajustes

□ Produção

□ Pedidos

□ Financeiro

□ Dashboard

□ Relatórios

□ Backup

□ Exportação

□ Configurações

Ao final:

Sistema totalmente operacional.

---

# Fase 3 — Experiência do Usuário

Objetivo:

Refinar toda experiência.

Inclui:

- microinterações

- loading

- estados vazios

- animações

- mensagens

- onboarding

- feedback visual

- responsividade

- acessibilidade

Meta:

Produto parecer premium.

---

# Fase 4 — Inteligência

Objetivo:

Adicionar IA ao sistema.

Possibilidades:

- previsão de estoque

- previsão financeira

- sugestão de compras

- previsão de vendas

- desperdício

- margem

- precificação

- alertas inteligentes

Nenhuma IA deverá alterar regras de negócio.

IA apenas recomenda.

---

# Fase 5 — Integrações

Objetivo:

Conectar o HorizonGest ao mundo externo.

Exemplos:

PIX

WhatsApp

iFood

Delivery

Mercado Pago

Stone

PagSeguro

NFC-e

SAT

Balanças

Impressoras

ERP

API pública

Webhook

---

# Fase 6 — Comercial

Objetivo:

Transformar em SaaS.

Inclui:

Planos

Assinaturas

Cobrança

Marketplace

Trial

Licenciamento

White Label

Revendedores

Afiliados

---

# Fase 7 — Escalabilidade

Objetivo:

Preparar milhares de empresas.

Inclui:

Redis

PostgreSQL Cluster

Filas

Workers

Horizontal Scaling

CDN

Observabilidade

Tracing

Monitoramento

---

# Ordem Obrigatória

Nenhuma fase futura pode comprometer uma fase anterior.

Nunca sacrificar:

Arquitetura

Engenharia

Segurança

Performance

Multi-tenant

para entregar funcionalidades.

Fim da Parte 26.


# PARTE 27

# Checklist Oficial de Desenvolvimento

Todo módulo novo deve seguir exatamente este checklist.

Nenhuma etapa pode ser ignorada.

---

# 1. Modelagem

□ Entidade criada

□ Domain criado

□ Interfaces criadas

□ CompanyID quando necessário

□ Soft Delete quando necessário

□ Auditoria quando necessário

---

# 2. Banco

□ Migration criada

□ Índices criados

□ Constraints criadas

□ Chaves estrangeiras

□ UNIQUE quando necessário

---

# 3. Repository

□ Interface

□ Implementação GORM

□ Todos os filtros por tenant

□ Context propagado

□ Sem regra de negócio

---

# 4. Service

□ Interface

□ Implementação

□ Toda regra de negócio

□ Validações

□ Erros padronizados

□ Transações quando necessário

---

# 5. Handler

□ Endpoint

□ Request

□ Response

□ Status HTTP corretos

□ Validação

□ Swagger atualizado

---

# 6. Frontend

□ Página

□ Componentes

□ Estados

□ Loading

□ Empty State

□ Error State

□ Toasts

□ Responsivo

---

# 7. Permissões

□ RBAC

□ Roles

□ Permissões

□ Menu

□ Rotas protegidas

---

# 8. Testes

□ Unitários

□ Integração

□ Fluxo principal

---

# 9. Documentação

□ API

□ AI_CONTEXT

□ Manual

□ DECISIONS (se necessário)

---

# 10. Revisão

Antes do merge verificar:

□ Arquitetura preservada

□ Nenhuma duplicação

□ Nenhum acoplamento novo

□ Nenhuma regra no frontend

□ Nenhuma regra no handler

□ Nenhuma consulta sem tenant

□ Nenhum branding hardcoded

□ Código limpo

□ Build funcionando

□ Testes passando

---

## Definition of Done

Um módulo só é considerado pronto quando:

- Funciona
- Está documentado
- Está testado
- Está protegido
- Está desacoplado
- Está dentro da arquitetura
- Está aprovado na revisão

Fim da Parte 27.


# PARTE 28

# Roadmap Oficial do Produto

Este roadmap representa a evolução planejada do HorizonGest.

Mudanças somente mediante decisão arquitetural registrada.

---

# FASE 1 — Fundação (Concluída)

Status:
✅ CONCLUÍDA

Incluiu:

- Arquitetura
- Multi-tenant
- Plataforma
- RBAC
- Branding
- White Label
- Documentação
- Manuais
- Auditoria
- Foundation Closed

---

# FASE 2 — MVP Comercial

Status:
Em andamento

Objetivo:

Permitir operação real de restaurantes.

Inclui:

## Cadastros

- Empresas
- Usuários
- Produtos
- Categorias
- Ingredientes
- Compras
- Estoque
- Clientes

## Operação

- Pedidos
- Produção
- Consumo de estoque
- Cancelamentos
- Histórico

## Financeiro

- Receitas
- Despesas
- Fluxo de caixa
- Contas

## Dashboard

- KPIs
- Resumos
- Alertas

---

# FASE 3 — Estabilização

Objetivo:

Eliminar bugs.

Checklist:

- Fluxos completos
- Performance
- UX
- Logs
- Tratamento de erros
- Ajustes de telas
- Revisão completa

Nenhuma funcionalidade grande deve entrar aqui.

---

# FASE 4 — Inteligência

Objetivo:

Adicionar inteligência operacional.

Inclui:

- Sugestões automáticas
- Indicadores
- Alertas inteligentes
- IA
- Recomendações

---

# FASE 5 — Integrações

Objetivo:

Conectar o HorizonGest ao ecossistema.

Possíveis integrações:

- WhatsApp
- PIX
- Bancos
- ERP
- Delivery
- Mercado Pago
- Stone
- PagSeguro
- iFood
- APIs fiscais

---

# FASE 6 — Plataforma

Objetivo:

Transformar HorizonGest em plataforma.

Inclui:

- Marketplace
- Plugins
- Extensões
- API Pública
- SDK

---

# FASE 7 — Escala

Objetivo:

Preparar milhares de empresas.

Inclui:

- PostgreSQL
- Redis
- Filas
- CDN
- Balanceamento
- Observabilidade
- Monitoramento
- Cluster

---

# Prioridade Oficial

Sempre seguir esta ordem:

1.
Corrigir Bugs

2.
Fluxo quebrado

3.
Usabilidade

4.
Performance

5.
Novas funcionalidades

Nunca inverter esta prioridade.

---

# Filosofia do Produto

Primeiro:

Funcionar.

Depois:

Ficar bonito.

Depois:

Escalar.

Nunca o contrário.

---

# Critério para novas funcionalidades

Uma funcionalidade só entra quando:

✓ resolve um problema real

✓ respeita a arquitetura

✓ não gera acoplamento

✓ possui documentação

✓ possui testes

✓ possui roadmap

---

Fim da Parte 28.


# PARTE 29

# Constituição Oficial do HorizonGest

Versão 1.0

---

# Preâmbulo

O HorizonGest não é apenas um software.

É uma plataforma construída para evoluir durante muitos anos.

Todas as decisões deste projeto devem preservar sua capacidade de crescer sem perder simplicidade.

A arquitetura existe para proteger o produto.

Não para limitar sua evolução.

---

# Artigo 1

## A Arquitetura é Soberana

Nenhuma funcionalidade possui autoridade para alterar a arquitetura.

Se uma funcionalidade exigir quebrar a arquitetura,

a funcionalidade deverá ser redesenhada.

Nunca a arquitetura.

---

# Artigo 2

## Simplicidade

A solução mais simples que resolve completamente o problema sempre será preferida.

Complexidade somente quando absolutamente necessária.

---

# Artigo 3

## Clareza

Código deve ser escrito para pessoas.

Computadores apenas executam.

Desenvolvedores precisam compreender.

---

# Artigo 4

## Multi-Tenant

Toda empresa deve existir como se fosse a única utilizando o sistema.

Nenhuma empresa poderá descobrir informações de outra.

Nem direta.

Nem indiretamente.

---

# Artigo 5

## Engenharia

O HorizonGest possui engenharia definida.

Ela deve ser respeitada.

Não existem atalhos.

---

# Artigo 6

## Regras de Negócio

Toda regra pertence ao Backend.

Mais especificamente ao Service.

Sempre.

---

# Artigo 7

## Frontend

O Frontend existe para proporcionar experiência.

Nunca para tomar decisões de negócio.

---

# Artigo 8

## Banco de Dados

O banco é uma implementação.

Não uma regra de negócio.

Toda lógica pertence ao domínio.

---

# Artigo 9

## White Label

O sistema pertence ao cliente.

Nunca ao desenvolvedor.

Toda identidade visual deve ser configurável.

---

# Artigo 10

## Documentação

Se não está documentado,

não faz parte oficialmente do projeto.

---

# Artigo 11

## Testes

Código sem testes é código parcialmente confiável.

Sempre que possível,

novas regras devem possuir testes.

---

# Artigo 12

## Segurança

Segurança nunca é funcionalidade.

É requisito.

---

# Artigo 13

## Performance

Performance é consequência de boa engenharia.

Nunca objetivo isolado.

---

# Artigo 14

## Dívida Técnica

Pode existir.

Jamais escondida.

Sempre registrada.

Sempre conhecida.

---

# Artigo 15

## Inteligência Artificial

Toda IA utilizada no projeto deve obedecer integralmente este Manual.

Nenhuma IA possui autonomia para alterar arquitetura ou engenharia.

---

# Artigo 16

## Evolução

O HorizonGest deve conseguir crescer indefinidamente.

Novos módulos devem parecer que sempre fizeram parte do sistema.

---

# Artigo 17

## Qualidade

"Funcionar"

não significa

"estar pronto".

Estar pronto significa:

- funcionar

- estar documentado

- estar testado

- respeitar arquitetura

- respeitar engenharia

- respeitar segurança

---

# Artigo 18

## Filosofia

Preferimos:

clareza

à esperteza.

Consistência

à criatividade.

Arquitetura

à velocidade.

Manutenção

à complexidade.

Longo prazo

ao improviso.

---

# Artigo 19

## Missão

Construir um ERP moderno,

escalável,

multi-tenant,

white-label,

capaz de evoluir continuamente

sem perder organização.

---

# Artigo 20

## Regra Suprema

Sempre que existir conflito entre:

Funcionalidade

e

Arquitetura

vence a Arquitetura.

Sempre.

---

# Encerramento

Este Manual representa a referência oficial do HorizonGest.

Qualquer decisão técnica futura deverá estar alinhada com este documento.

Mudanças somente poderão ocorrer mediante decisão arquitetural formal registrada no DECISIONS.md.

Este documento deverá acompanhar toda a vida útil do projeto.

---

Fim do Manual Mestre HorizonGest.

Versão 1.0

Status:

**OFICIAL**



