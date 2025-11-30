# Architecture Overview
## BRIN GenAI Workshop - Phase 1.2

**Duration**: 20 minutes
**Format**: 10 min explanation + 5 min Q&A + 5 min setup check

---

## 🎯 Objective

Understand system components, integration points, and technology stack rationale

---

## 📊 Slide 1: System Architecture Overview

**Complete System - 3 Layers**

```mermaid
graph TB
    subgraph "User Layer"
        U[WhatsApp User]
    end

    subgraph "Integration Layer"
        WM[Whatmeow<br/>WA Handler]
        DB[(PostgreSQL<br/>User & History)]
    end

    subgraph "Automation Layer"
        WF1[Webhook Trigger]
        WF2[User Context]
        WF3[Intent Detection]
        WF4[Routing Logic]
        WF5[Response Builder]
    end

    subgraph "Intelligence Layer"
        OR[OpenRouter<br/>LLM Gateway]
        RAG[RAG System]
        VDB[(Vector DB<br/>pgvector)]
    end

    U -->|Message| WM
    WM -->|Check User| DB
    WM -->|Trigger| WF1
    WF1 --> WF2
    WF2 --> WF3
    WF3 --> WF4
    WF4 -->|Simple| OR
    WF4 -->|Knowledge| RAG
    RAG --> VDB
    RAG --> OR
    OR --> WF5
    WF5 --> WM
    WM -->|Reply| U

    style U fill:#e3f2fd
    style WM fill:#fff3e0
    style WF4 fill:#f3e5f5
    style OR fill:#c8e6c9
    style RAG fill:#fce4ec
```

---

## 📊 Slide 2: Data Flow Sequence

**Example: User Query with RAG**

```mermaid
sequenceDiagram
    participant U as User
    participant BE as Golang Backend
    participant DB as PostgreSQL
    participant N8N as N8N Workflow
    participant RAG as RAG System
    participant LLM as OpenRouter

    U->>BE: "How do I submit proposal?"
    BE->>DB: Get user context
    DB-->>BE: User data
    BE->>N8N: Trigger + context
    N8N->>N8N: Detect intent: Knowledge Query
    N8N->>RAG: Search similar docs
    RAG-->>N8N: Top 3 relevant docs
    N8N->>LLM: Generate (prompt + context)
    LLM-->>N8N: Custom response
    N8N->>DB: Log interaction
    N8N-->>BE: Response
    BE->>U: Send reply
```

**Key Points**:
- User context enrichment
- Intent-based routing
- RAG for knowledge retrieval
- LLM for generation
- Database logging

---

## 📊 Slide 3: Component Breakdown

**What We'll Build Today**

```mermaid
graph TB
    subgraph "Pre-Built"
        C0[WhatsApp Backend<br/>Golang + Whatmeow]
    end

    subgraph "Hands-On Components"
        C1[Component 1<br/>N8N Workflow<br/>Automation Engine]
        C2[Component 2<br/>Prompt Templates<br/>LLM Instructions]
        C3[Component 3<br/>RAG Pipeline<br/>Knowledge Retrieval]
        C4[Component 4<br/>Integration<br/>End-to-End Flow]
    end

    C0 -.->|Provided| C1
    C1 -->|Module 1| C2
    C2 -->|Module 2| C3
    C3 -->|Module 3| C4

    style C0 fill:#e0e0e0
    style C1 fill:#e3f2fd
    style C2 fill:#f3e5f5
    style C3 fill:#fff3e0
    style C4 fill:#c8e6c9
```

**Focus Areas**:
1. ✅ N8N Workflow (Module 1) - Understand pre-built template
2. ✅ Prompt Engineering (Module 2) - Theory and best practices
3. ✅ RAG Implementation (Module 3) - Use pre-built LangChain workflow
4. ✅ Integration (Phase 3) - Connect all components

---

## 📊 Slide 4: Technology Stack

**Why These Technologies?**

| Component | Technology | Reason |
|-----------|-----------|---------|
| **Messaging** | Whatmeow | Open-source, reliable WhatsApp API |
| **Workflow** | N8N | Visual automation, easy debugging |
| **Database** | PostgreSQL + pgvector | Unified DB for data + vectors |
| **LLM Gateway** | OpenRouter | Multi-model access, cost optimization |
| **Embeddings** | OpenAI | High-quality vector representations |
| **Container** | Docker | Consistent environment, easy deployment |

**All Components Run Locally** - No cloud dependencies (except LLM APIs)

---

## 📊 Slide 5: Why N8N over Custom Code?

**Decision Rationale**

### ✅ Benefits of N8N
- **Visual Workflow**: See data flow in real-time
- **Rapid Development**: Build in hours vs days
- **Easy Debugging**: Step through execution logs
- **No Code Changes**: Modify workflows without deployment
- **Production Ready**: Used by thousands of companies
- **LangChain Integration**: Built-in RAG nodes

### ⚠️ Trade-offs
- Additional infrastructure dependency
- Less control over low-level implementation

### 💡 When to Use Custom Code Instead
- Need millisecond-level performance
- Complex custom logic
- Already have existing codebase

**For this workshop**: N8N is perfect for learning and prototyping

---

## 📊 Slide 6: Why OpenRouter?

**Multi-Model LLM Gateway**

```mermaid
graph LR
    App[Your Application]
    OR[OpenRouter<br/>Single API]

    OR --> GPT[GPT-4o<br/>Creative Tasks]
    OR --> Claude[Claude 3.5<br/>Reasoning]
    OR --> Llama[Llama 3.1<br/>Fast/Cheap]
    OR --> Other[20+ Models]

    App -->|One API Key| OR

    style App fill:#e3f2fd
    style OR fill:#fff3e0
    style GPT fill:#c8e6c9
    style Claude fill:#f3e5f5
    style Llama fill:#fce4ec
```

**Advantages**:
- ✅ Single API for multiple models
- ✅ Automatic fallback if one provider down
- ✅ Cost optimization (route to cheapest model)
- ✅ Easy A/B testing between models
- ✅ No vendor lock-in

---

## 📊 Slide 7: Why RAG with pgvector?

**PostgreSQL + pgvector vs Dedicated Vector DB**

### ✅ Why pgvector?
- **Unified Database**: User data + vectors in same DB
- **Simplified Architecture**: One less service to manage
- **Zero Additional Cost**: No separate subscription
- **Familiar Technology**: BRIN team knows PostgreSQL
- **Docker Ready**: Single container with pgvector

### 🔍 When to Use Dedicated Vector DB (Pinecone/Weaviate)?
- Millions of vectors (100M+)
- Sub-10ms query requirements
- Complex vector operations

**For CS automation**: pgvector is perfect (< 100K vectors)

---

## 📊 Slide 8: Docker-Based Development

**One-Command Setup Philosophy**

```bash
# Clone repository
git clone <repo-url>
cd workshop-brin

# Start entire stack
docker-compose up -d

# Verify services
docker-compose ps
```

**All Services in Containers**:
```mermaid
graph TB
    subgraph "Docker Compose"
        DC[docker-compose.yaml]

        DC --> BE[backend-wa<br/>Golang + Whatmeow<br/>Port: 8080]
        DC --> N8N[n8n<br/>Workflow Engine<br/>Port: 5678]
        DC --> PG[postgres<br/>PostgreSQL + pgvector<br/>Port: 5432]
    end

    style DC fill:#e3f2fd
    style BE fill:#fff3e0
    style N8N fill:#c8e6c9
    style PG fill:#f3e5f5
```

**Benefits**:
- ✅ Zero OS dependencies (only Docker)
- ✅ Identical environment across Windows/Mac/Linux
- ✅ Easy teardown: `docker-compose down`
- ✅ Take home and continue learning

---

## 📊 Slide 9: Project Structure

**Repository Organization**

```
workshop-brin/
├── backend/              # Golang WhatsApp service (pre-built)
│   ├── main.go
│   ├── handlers/
│   └── Dockerfile
├── docs/                 # Documentation
│   └── setup-guide.md
├── n8n-workflows/       # Importable N8N workflows
│   ├── 01-basic-llm.json
│   ├── 02-rag-ingestion.json
│   └── 03-rag-query.json
├── scripts/             # Database initialization
│   ├── init-db.sql
│   ├── setup-pgvector.sql
│   └── seed-data.sql
├── knowledge-base/      # Sample documents
│   ├── sample-faqs.txt
│   └── policies.txt
├── docker-compose.yaml  # Complete stack
└── .env.example        # Environment variables
```

**Everything Pre-Configured** - Just import and run!

---

## 📊 Slide 10: Workshop Timeline Recap

**What Happens Next (90 minutes)**

```mermaid
gantt
    title Workshop Phase 2 & 3
    dateFormat mm
    axisFormat %M min

    section Module 1
    N8N + LLM (Template)    :00, 30m

    section Module 2
    Prompt Engineering (Theory)   :30, 25m

    section Module 3
    RAG Implementation (Template)  :55, 35m

    section Phase 3
    Integration & Demo             :90, 30m
```

**Approach**:
- Module 1: Import pre-built N8N workflow, understand concepts
- Module 2: Learn prompt engineering (no hands-on, theory + examples)
- Module 3: Import RAG workflow with LangChain, test with queries
- Phase 3: Connect all components, test end-to-end

---

## 📊 Slide 11: Q&A

**Common Questions**

**Q**: Do I need to know machine learning?
**A**: No! We use pre-trained models via APIs

**Q**: Can this work with Telegram/Slack?
**A**: Yes, architecture is platform-agnostic. Change Whatmeow to other integrations.

**Q**: What about costs in production?
**A**: We'll cover cost optimization in Module 2 (model selection)

**Q**: Is this production-ready?
**A**: Yes! N8N + OpenRouter + pgvector is used in production by many companies

**Q**: Can I use local LLMs instead of OpenRouter?
**A**: Yes! Ollama integration is possible (not covered today)

---

## 📊 Slide 12: Pre-Workshop Setup Check

**Verify Your Environment** (5 minutes)

### ✅ Checklist:

1. **Docker Running**
   ```bash
   docker ps  # Should return without error
   ```

2. **Repository Cloned**
   ```bash
   cd workshop-brin
   ls  # Should see docker-compose.yaml
   ```

3. **Environment Variables**
   ```bash
   cp .env.example .env
   # Add your OpenRouter API key
   ```

4. **Start Services**
   ```bash
   docker-compose up -d
   docker-compose ps  # All services should be "Up"
   ```

5. **Verify N8N**
   - Open browser: `http://localhost:5678`
   - Should see N8N login screen

**Teaching Assistants**: Please help participants who encounter issues!

---

## 📊 Slide 13: Transition to Module 1

**Up Next: N8N Workflow + LLM Integration**

**What You'll Do** (30 minutes):
- Import pre-built N8N workflow
- Understand webhook triggers
- Learn OpenRouter API integration
- Observe workflow execution
- See LLM request/response patterns

**Get Ready**:
- Open N8N: `http://localhost:5678`
- Have your OpenRouter API key ready
- Follow along with instructor

**Let's Build!** 🚀

---

## 🎓 Instructor Notes

**Timing Guidelines**:
- Slides 1-3: 5 minutes (architecture overview)
- Slides 4-7: 5 minutes (technology decisions)
- Slides 8-10: 3 minutes (Docker setup and project structure)
- Slide 11: 5 minutes (Q&A - be flexible)
- Slides 12-13: 2 minutes (setup verification and transition)

**Key Points to Emphasize**:
1. **Architecture is modular** - each component can be swapped
2. **N8N makes automation accessible** - visual debugging is powerful
3. **Docker ensures consistency** - no "works on my machine" problems
4. **OpenRouter provides flexibility** - switch models easily
5. **pgvector is practical** - no need for complex vector DB setup

**Demo During Q&A**:
- Show N8N interface briefly
- Show docker-compose.yaml
- Show how to check Docker logs

**Watch For**:
- Participants struggling with Docker setup
- Missing .env file configuration
- Firewall blocking localhost ports
- Teaching assistants should help during setup check

**Transition Smoothly**:
- "Now that you understand the architecture, let's build it!"
- "Module 1 starts with N8N - the brain of our automation"
- Ensure all participants have N8N running before proceeding