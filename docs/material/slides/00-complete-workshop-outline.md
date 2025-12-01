# BRIN GenAI Workshop - Complete Outline
**AI-Powered Customer Service Automation**

**Duration**: 2.5 hours | **Target**: Mid-Level Developers | **Focus**: 65% Hands-on / 35% Theory

---

## Phase 1: Foundation (30 minutes)

Set context and verify technical setup for the workshop. Participants understand the GenAI paradigm shift, system architecture, and ensure their development environment is ready.

### 01. Opening: AI → GenAI → Practical Impact (10 min)
**[File: 01-opening-genai-intro.md](01-opening-genai-intro.md)**

- Slide 1: Title Slide - Workshop Overview
- Slide 2: Traditional AI vs Generative AI
- Slide 3: The GenAI Paradigm Shift
- Slide 4: Real-World Impact
- Slide 5: What is Automation Workflow?
- Slide 6: What We'll Build Today
- Slide 7: Workshop Flow
- Slide 8: Key Takeaways
- Slide 9: Transition to Next Section
- Live Demo: Working WhatsApp Bot (FAQ, RAG Query, Complaint)

### 02. Architecture Overview (20 min)
**[File: 02-architecture-overview.md](02-architecture-overview.md)**

- Slide 1: System Architecture Overview (3 Layers)
- Slide 2: Data Flow Sequence
- Slide 3: Component Breakdown
- Slide 4: Technology Stack
- Slide 5: Why N8N over Custom Code?
- Slide 6: Why OpenRouter?
- Slide 7: Why RAG with pgvector?
- Slide 8: Docker-Based Development
- Slide 9: Project Structure
- Slide 10: Workshop Timeline Recap
- Slide 11: Q&A
- Slide 12: Pre-Workshop Setup Check
- Slide 13: Transition to Module 1

---

## Phase 2: Core Implementation (90 minutes)

Build the three core components: N8N workflow automation with LLM, learn prompt engineering principles, and implement RAG system with vector embeddings.

### 03. Module 1: N8N Workflow + LLM Integration (30 min)
**[File: 03-module1-n8n-llm.md](03-module1-n8n-llm.md)**

- Slide 1: Module Overview
- Slide 2: N8N Workflow Architecture (6 Nodes)
- Slide 3: Node 1 - Webhook Trigger
- Slide 4: LLM API Integration (OpenRouter)
- Slide 5: Process Response
- Slide 6: N8N Key Concepts
- Slide 7: Testing the Workflow
- Slide 8: Hands-On Activity (Import, Configure, Test)
- Slide 9: Key Takeaways
- Slide 10: Transition to Module 2

### 04. Module 2: Prompt Engineering (25 min)
**[File: 04-module2-prompt-engineering.md](04-module2-prompt-engineering.md)**

- Slide 1: Module Overview
- Slide 2: Prompt Engineering Framework (4-Part Structure)
- Slide 3: Template Structure
- Slide 4: Scenario 1 - FAQ Handler
- Slide 5: Scenario 2 - Complaint Handler
- Slide 6: Scenario 3 - Escalation Detector
- Slide 7: Advanced Technique 1 - Chain-of-Thought (CoT)
- Slide 8: Advanced Technique 2 - Few-Shot Examples
- Slide 9: Advanced Technique 3 - ReAct (Reasoning + Action)
- Slide 10: Parameter Tuning Guide
- Slide 11: Model Selection Decision Tree
- Slide 12: Common FAQ
- Slide 13: Key Takeaways
- Slide 14: Transition to Module 3

### 05. Module 3: RAG Implementation (35 min)
**[File: 05-module3-rag-implementation.md](05-module3-rag-implementation.md)**

- Slide 1: What is RAG?
- Slide 2: Why RAG?
- Slide 3: How RAG Works (Two Phases)
- Slide 4: Vector Embeddings Explained
- Slide 5: N8N LangChain Integration
- Slide 6: RAG Ingestion Workflow
- Slide 7: Text Chunking Strategy
- Slide 8: RAG Query Workflow
- Slide 9: RAG Prompt Template
- Slide 10: Similarity Search in Action
- Slide 11: Hands-On Activity - Part A (Document Ingestion)
- Slide 12: Hands-On Activity - Part B (RAG Query)
- Slide 13: Testing & Validation
- Slide 14: RAG Best Practices
- Slide 15: Common RAG Pitfalls
- Slide 16: Key Takeaways
- Slide 17: Transition to Phase 3

---

## Phase 3: Integration (30 minutes)

Connect all components into a working end-to-end system with intent detection, routing logic, and production considerations.

### 06. End-to-End Integration & Demo (20 min)
**[File: 06-integration-demo.md](06-integration-demo.md)**

- Slide 1: What We've Built So Far
- Slide 2: Integration Architecture
- Slide 3: Intent Detection Logic
- Slide 4: Intent Detection Implementation
- Slide 5: Routing Workflow
- Slide 6: Hands-On Integration (15 min activity)
- Slide 7: Test Scenarios
- Slide 8: Troubleshooting Common Issues
- Slide 9: Live Demo - Golang Backend Integration
- Slide 10: Performance Metrics
- Slide 11: Production Deployment Considerations
- Slide 12: Key Takeaways
- Slide 13: What's Next?
- Slide 14: Transition to Closing

### 07. Closing & Next Steps (10 min)
**[File: 07-closing-next-steps.md](07-closing-next-steps.md)**

- Slide 1: Workshop Journey Recap
- Slide 2: Key Learnings
- Slide 3: What You Can Build Now
- Slide 4: Your Complete System
- Slide 5: Next Steps - Immediate Actions
- Slide 6: Next Steps - Advanced Features
- Slide 7: Production Deployment Checklist
- Slide 8: Learning Resources
- Slide 9: Cost Optimization Tips
- Slide 10: Common Pitfalls to Avoid
- Slide 11: Support & Community
- Slide 12: Success Stories to Inspire
- Slide 13: Final Thoughts
- Slide 14: Q&A Session
- Slide 15: Thank You!

---

## Summary

**Total Duration**: 150 minutes (2.5 hours)

**Time Breakdown**:
- Phase 1 (Foundation): 30 min (20%)
- Phase 2 (Core Implementation): 90 min (60%)
- Phase 3 (Integration & Closing): 30 min (20%)

**Key Technologies**:
N8N, OpenRouter, PostgreSQL + pgvector, LangChain, Whatmeow, Docker

**Deliverables**:
4 Working N8N Workflows, Complete Docker Stack, Sample Knowledge Base, Documentation