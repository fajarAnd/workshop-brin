# Claude Code Project Rules & Conventions

## 📋 **Planning & Task Management Standards**

### **Spec Then Design Methodology**
Use the "STD" prefix in prompts to trigger formal specification-first development:

### **Specification Document Requirements**
- **Document location**: `docs/specs/design` directory
- **File naming format**: `YYYYMMDDHHMM-spec-name.md` (e.g., `202508031830-news-scraping-endpoint.md`)
- **Template structure**: Use standardized template below
- **Review process**: NEVER implement code directly - always present specification to user for review and discussion first
- **Minimize code in documentation**: Only include code snippets when absolutely necessary and cannot be represented by diagrams
- **Visualization-first approach**: Use Mermaid JS diagrams for workflows, architecture, data flow, and logic visualization
- **Diagram standards**: All diagrams must be clean, readable, and simple - avoid complexity
- **Decision log required**: All architectural and design decisions must be documented in Decision Log section

### **STC Planning Document Template** (Only when "stc" prompted)
```markdown
# [Component] - Technical Specification

**Document ID**: YYYYMMDDHHMM-component-name
**Author**: Claude Code
**Date**: YYYY-MM-DD
**Status**: Planning

## Table of Contents

- [Problem Statement](#problem-statement)
- [Current State Analysis](#current-state-analysis)
- [Core Requirements](#core-requirements)
- [Technical Design](#technical-design)
- [Implementation Plan](#implementation-plan)
  - [Phase 1: Core Implementation](#phase-1-core-implementation)
  - [Phase 2: Integration](#phase-2-integration)
- [Decision Log](#decision-log)
- [Verification Criteria](#verification-criteria)

## Problem Statement
What needs to be solved and why.

## Current State Analysis
Analysis of existing implementation and gaps.

## Core Requirements
- REQ-001: [Essential functionality]
- REQ-002: [Critical business rule]

## Technical Design
Architecture overview and component design.

**Use Mermaid JS diagrams** for:
- System architecture
- Data flow sequences
- Component relationships
- Workflow visualization
- State transitions

**Minimize code snippets** - only include when:
- Code structure is critical to understanding
- Cannot be represented visually
- Shows important implementation pattern

## Implementation Plan
### Phase 1: Opening & Foundation
- [ ] Intro: AI → GenAI → Practical Use Case
- [ ] General Architecture

### Phase 2: Prompt Engineering for CS Automation
- [ ] Basic Prompt Engineering
- [ ] Effective Prompting

## Decision Log
Document all architectural and design decisions:

| Decision | Rationale | Alternatives Considered |
|----------|-----------|------------------------|
| [Technology/approach chosen] | [Why this was selected] | [What else was considered] |
```

### **Task Management Workflow**
1. **STD Tasks**: 
   - Create specification document in `docs/specs/design` using required format
   - Present specification to user for review and discussion
   - Wait for user approval before implementation
   - NEVER implement before specification review
2. **Regular Tasks**: Use TodoWrite tool for multi-step tasks (3+ steps)
3. **Simple Tasks**: Direct implementation without formal planning
4. **Task Tracking**: 
   - Mark as `in_progress` when starting
   - Mark as `completed` when finished
   - Only one complex task `in_progress` at a time

---
*This file defines the coding standards and conventions for this project. Claude Code will follow these rules automatically.*






Revisi:
- Format Hands-on Laptop masing-masing. template project akan di provide ( akan di dijelaskan di bawah) 
- Open router API key Akan di sediakan
- untuk database jangan gunakan supabase, ganti saja menjadi postgres dengan menggunakan extension pg vector.
- N8N Self hosted di local
- Project template dengan structure di bawah
```
.
├── backend/           # Go Whatapp service (whatsmeow)
├── docs/             # Architecture documentation
├── scripts/          # Database initialization scripts
├── docker-compose.yaml
└── .env.example
```
- untuk backend akan dibuat manual oleh saya, jangan masukan ke spec.
- Project ini bisa di running  dengan eksekusi 1 perintah. gunakan "docker-compose up -d". usahakan minimal effort, tidak ada dependency (OS, external library etc) hanya cukup perlu docker
- docker-compose memiliki service: backend-wa, n8n, postgres, 

Bertahap, slide persentation dulu berbentuk markdown. setelah saya review baru step by step untuk detail dan comprehenshive doc tersebut.

