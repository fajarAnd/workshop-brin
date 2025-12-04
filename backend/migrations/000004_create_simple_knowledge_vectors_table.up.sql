CREATE TABLE public.simple_knowledge_vectors (
                                                 id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                                                 content TEXT NOT NULL,
                                                 embedding VECTOR(1536), -- OpenAI embedding dimension
                                                 metadata JSONB,
                                                 created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                                                 updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);