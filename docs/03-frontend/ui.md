# UI Documentation

**HorizonGest Platform - Frontend UI Guidelines**

---

## Design Language

### Overview

The HorizonGest design language was created to convey a professional, modern, and reliable experience. Inspired by successful SaaS products like Stripe, Linear, and Vercel, the system aims to abandon the academic CRUD appearance in favor of a refined commercial interface.

### Principles

**Lightness**
- Generous whitespace
- Subtle visual elements without excessive borders and boxes
- Clear hierarchy through spacing, not visual weight

**Speed**
- Smooth micro-interactions that respond immediately
- Instant visual feedback on all interactions
- Short animations (150-300ms) that don't hurt performance

**Organization**
- Consistent grid throughout the application
- Logical grouping of related information
- Intuitive navigation with clear context indicators

**Professionalism**
- Refined typography with strong visual hierarchy
- Sophisticated color palette
- Attention to detail in every pixel

**Minimalism**
- "Less is more" - remove everything non-essential
- Focus on content and primary actions
- Clean interface without unnecessary distractions

**Reliability**
- Clear and elegant loading states
- Informative and constructive error messages
- Visual feedback on all user actions

### Brand Personality

HorizonGest is:
- **Sophisticated**: Not a toy, it's business
- **Reliable**: Users can count on the system
- **Efficient**: Helps users work fast
- **Modern**: Uses current best design practices
- **Accessible**: Easy to use for any operator

### Tone of Voice

- **Direct**: Get to the point
- **Professional**: Use language appropriate to context
- **Friendly**: Be accessible, not robotic
- **Clear**: Avoid unnecessary technical jargon

---

## Layout

### Grid System

- **Base Grid:** 12-column grid
- **Gutter:** 24px
- **Container Max Width:** 1200px
- **Sidebar Width:** 280px (collapsed: 64px)

### Sidebar

**Structure:**
- Logo area at top
- Navigation groups
- User profile at bottom
- Collapse/expand toggle

**Navigation Groups:**
- Main navigation
- Settings
- Help

### Content Areas

**Cards:**
- Subtle borders
- Generous padding (24px)
- Clear hierarchy
- Hover states

**Tables:**
- Clean rows with subtle borders
- Sortable headers
- Pagination
- Horizontal scroll on mobile

---

## Components

### Buttons

**Primary Button:**
- Background: Primary color
- Text: White
- Padding: 10px 20px
- Border radius: 6px
- Hover: Darken primary color

**Secondary Button:**
- Background: Transparent
- Text: Primary color
- Border: 1px solid primary color
- Padding: 10px 20px
- Border radius: 6px
- Hover: Light background

**Ghost Button:**
- Background: Transparent
- Text: Text color
- Border: None
- Padding: 10px 20px
- Hover: Light background

### Inputs

**Text Input:**
- Border: 1px solid #e5e7eb
- Border radius: 6px
- Padding: 10px 12px
- Focus: Primary color border
- Error: Red border

**Select:**
- Same as text input
- Custom dropdown arrow

**Checkbox:**
- Custom styled checkbox
- Primary color when checked

### Cards

**Structure:**
- Header (optional)
- Body
- Footer (optional)

**Styling:**
- Border: 1px solid #e5e7eb
- Border radius: 8px
- Padding: 24px
- Background: White
- Shadow: Subtle on hover

### Modals

**Structure:**
- Overlay (dark, semi-transparent)
- Modal container
- Header (title + close button)
- Body
- Footer (actions)

**Styling:**
- Max width: 500px
- Border radius: 12px
- Animation: Fade in/slide up

---

## Colors

### Primary Palette

- **Primary:** #0f172a (dark blue)
- **Secondary:** #6366f1 (indigo)
- **Accent:** #3b82f6 (blue)

### Neutral Palette

- **Background:** #ffffff
- **Surface:** #f9fafb
- **Border:** #e5e7eb
- **Text Primary:** #111827
- **Text Secondary:** #6b7280
- **Text Tertiary:** #9ca3af

### Semantic Colors

- **Success:** #10b981 (green)
- **Warning:** #f59e0b (amber)
- **Error:** #ef4444 (red)
- **Info:** #3b82f6 (blue)

### Dark Mode

- **Background:** #0f172a
- **Surface:** #1e293b
- **Border:** #334155
- **Text Primary:** #f8fafc
- **Text Secondary:** #cbd5e1

---

## Typography

### Font Family

- **Primary:** Inter (system-ui fallback)
- **Monospace:** JetBrains Mono (monospace fallback)

### Font Sizes

- **Display:** 48px (H1)
- **Heading:** 32px (H2)
- **Title:** 24px (H3)
- **Subtitle:** 18px (H4)
- **Body:** 16px (p)
- **Small:** 14px (small)
- **Tiny:** 12px (label)

### Font Weights

- **Regular:** 400
- **Medium:** 500
- **Semibold:** 600
- **Bold:** 700

### Line Heights

- **Tight:** 1.25 (headings)
- **Normal:** 1.5 (body)
- **Relaxed:** 1.75 (long text)

---

## Spacing

### Scale

- **0:** 0px
- **1:** 4px
- **2:** 8px
- **3:** 12px
- **4:** 16px
- **5:** 20px
- **6:** 24px
- **8:** 32px
- **10:** 40px
- **12:** 48px
- **16:** 64px

### Usage

- **Component Padding:** 24px (6)
- **Section Margin:** 48px (12)
- **Element Gap:** 16px (4)
- **Text Gap:** 8px (2)

---

## Responsive

### Breakpoints

- **Mobile:** < 640px
- **Tablet:** 640px - 1024px
- **Desktop:** 1024px - 1280px
- **Large:** > 1280px

### Mobile Adaptations

- Sidebar: Collapsed by default
- Tables: Horizontal scroll
- Cards: Full width
- Modals: Full screen

---

## Icons

### Icon Library

- **Library:** Lucide Svelte
- **Size:** 16px, 20px, 24px
- **Color:** Current text color

### Usage

- **Navigation:** 20px
- **Actions:** 16px
- **Headers:** 24px

---

## What NOT to Do

### Avoid

- ❌ Excessive colors and gradients
- ❌ Heavy borders and excessive boxes
- ❌ Long and distracting animations
- ❌ Emojis as main icons
- ❌ Generic template-like layouts
- ❌ All caps text (except discrete labels)
- ❌ Heavy and unnatural shadows
- ❌ Decorative or illegible fonts
- ❌ Inconsistent spacing
- ❌ Generic loading states ("Loading...")

### Copy vs. Inspire

- ✅ **Inspire:** Understand principles and adapt
- ❌ **Copy:** Reproduce identical layouts
- ✅ **Mix ideas:** Combine the best of each reference
- ❌ **Clone:** Make a replica of an existing product
- ✅ **Create unique identity:** Develop unique language
- ❌ **Reproduce:** Make exact copy of any reference

---

## Implementation Guidelines

### Performance

- Keep bundle size small
- Avoid unnecessary heavy libraries
- Leverage existing components
- Optimize re-renders
- Use native CSS when possible

### Accessibility

- Minimum contrast of 4.5:1 for text
- Functional keyboard navigation
- Descriptive labels on all inputs
- Visible focus states
- Alt text on images

### Maintainability

- Modular and reusable components
- CSS variables for design tokens
- Consistent naming
- Updated documentation
- Clean and organized code

---

## Success Metrics

The success of the design will be measured by:

1. **User Perception:** "This looks like professional software"
2. **Operational Efficiency:** Fewer clicks to complete tasks
3. **Visual Satisfaction:** Pleasant interface to use
4. **Performance:** Load time and interactivity
5. **Consistency:** Visual coherence throughout the application

---

## Evolution

The design language is alive and should evolve with:

- User feedback
- New business needs
- Industry best practices
- Emerging technologies
- Continuous learning

---

**Last Updated:** Fase 2 - Documentation & Knowledge Base
