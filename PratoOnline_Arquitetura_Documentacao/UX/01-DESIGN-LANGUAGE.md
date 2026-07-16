# Design Language - PratoOnline

## Visão Geral

O design language do PratoOnline foi criado para transmitir uma experiência profissional, moderna e confiável. Inspirado em produtos SaaS de sucesso como Stripe, Linear e Vercel, o sistema busca abandonar a aparência de CRUD acadêmico em favor de uma interface comercial refinada.

## Filosofia

### Princípios Fundamentais

**Leveza**
- Interface respirável com espaços em branco generosos
- Elementos visuais sutis sem excesso de bordas e caixas
- Hierarquia clara através de espaçamento, não de peso visual

**Rapidez**
- Microinterações suaves que respondem imediatamente às ações
- Feedback visual instantâneo em todas as interações
- Animações curtas (150-300ms) que não prejudicam a performance

**Organização**
- Grid consistente em toda a aplicação
- Agrupamento lógico de informações relacionadas
- Navegação intuitiva com indicadores claros de contexto

**Profissionalismo**
- Tipografia refinada com hierarquia visual forte
- Paleta de cores sóbria e sofisticada
- Atenção aos detalhes em cada pixel

**Minimalismo**
- "Menos é mais" - remover tudo o que não é essencial
- Foco no conteúdo e nas ações principais
- Interface limpa sem distrações desnecessárias

**Confiabilidade**
- Estados de loading claros e elegantes
- Mensagens de erro informativas e construtivas
- Feedback visual em todas as ações do usuário

## Identidade Visual

### Personalidade da Marca

O PratoOnline é:
- **Sofisticado**: Não é brincadeira, é negócio
- **Confiável**: O usuário pode contar com o sistema
- **Eficiente**: Ajuda o usuário a trabalhar rápido
- **Moderno**: Usa as melhores práticas atuais de design
- **Acessível**: Fácil de usar para qualquer operador

### Tom de Voz

- **Direto**: Vá direto ao ponto
- **Profissional**: Use linguagem adequada ao contexto
- **Amigável**: Seja acessível, não robótico
- **Claro**: Evite jargões técnicos desnecessários

## Referências de Design

### Inspirações Principais

**Autumn CRM Dashboard (Dribbble)**
- Layout horizontal com sidebar moderna
- Cards executivos com informações essenciais
- Uso inteligente de espaço em branco
- Hierarquia visual clara através de tipografia

**Artifact Dashboard**
- Componentes refinados e reutilizáveis
- Sidebar com ícones e grupos bem definidos
- Transições suaves entre estados
- Atenção aos detalhes de microinterações

**Stripe**
- Sistema de cores sofisticado
- Tipografia impecável
- Animações sutis e elegantes
- Dark mode bem implementado

**Linear**
- Minimalismo extremo
- Espaçamentos generosos
- Bordas sutis
- Performance excepcional

**Vercel**
- Interface limpa e profissional
- Componentes modulares
- Sistema de design consistente
- Acessibilidade em primeiro lugar

## O Que NÃO Fazer

### Evitar

- ❌ Excesso de cores e gradientes
- ❌ Bordas pesadas e caixas excessivas
- ❌ Animações longas e distrativas
- ❌ Emojis como ícones principais
- ❌ Layout que pareça template genérico
- ❌ Texto em todo maiúsculo (exceto labels discretas)
- ❌ Sombras pesadas e não naturais
- ❌ Fontes decorativas ou não legíveis
- ❌ Espaçamentos inconsistentes
- ❌ Estados de loading genéricos ("Carregando...")

### Copiar vs. Inspirar

- ✅ **Inspirar-se**: Entender os princípios e adaptar
- ❌ **Copiar**: Reproduzir layouts idênticos
- ✅ **Misturar ideias**: Combinar o melhor de cada referência
- ❌ **Clonar**: Fazer uma réplica de um produto existente
- ✅ **Criar identidade própria**: Desenvolver linguagem única
- ❌ **Reproduzir**: Fazer cópia fiel de qualquer referência

## Diretrizes de Implementação

### Performance

- Manter o bundle size pequeno
- Evitar bibliotecas pesadas desnecessárias
- Aproveitar componentes existentes
- Otimizar re-renderizações
- Usar CSS nativo sempre que possível

### Acessibilidade

- Contraste mínimo de 4.5:1 para texto
- Navegação por teclado funcional
- Labels descritivos em todos os inputs
- Estados de foco visíveis
- Texto alternativo em imagens

### Responsividade

- Mobile-first approach
- Breakpoints: 480px, 768px, 1024px, 1280px
- Sidebar adaptável
- Tabelas com scroll horizontal em mobile
- Touch targets mínimos de 44px

### Manutenibilidade

- Componentes modulares e reutilizáveis
- Variáveis CSS para design tokens
- Nomenclatura consistente
- Documentação atualizada
- Código limpo e organizado

## Métricas de Sucesso

O sucesso do novo design será medido por:

1. **Percepção do Usuário**: "Isso parece um software profissional"
2. **Eficiência Operacional**: Menos cliques para realizar tarefas
3. **Satisfação Visual**: Interface agradável de usar
4. **Performance**: Tempo de carregamento e interatividade
5. **Consistência**: Coerência visual em toda a aplicação

## Evolução

O design language é vivo e deve evoluir com:

- Feedback dos usuários
- Novas necessidades do negócio
- Melhores práticas da indústria
- Tecnologias emergentes
- Aprendizado contínuo

---

**Versão**: 1.0  
**Data**: 15/07/2026  
**Sprint**: 9 - Product Experience
