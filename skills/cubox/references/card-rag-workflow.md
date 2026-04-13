# RAG vs Keyword Search — Decision & Workflow

This reference defines how to choose between `card list --keyword` and `card rag --query`, and how to handle results progressively. **Follow these rules whenever the user asks you to find, search, or retrieve bookmarks from Cubox.**

## Step 1 — Choose the right search method


| User intent                                           | Method                     | Example                                                                                 |
| ----------------------------------------------------- | -------------------------- | --------------------------------------------------------------------------------------- |
| Exact term, title fragment, domain, known phrase      | `card list --keyword`      | "find my bookmarks from csdn.net", "search for 'React hooks'"                           |
| Conceptual question, topic exploration, fuzzy intent  | `card rag --query`         | "articles about building REST APIs with auth", "what did I save about LLM fine-tuning?" |
| Browse / filter by metadata (folder, tag, star, time) | `card list` (with filters) | "show starred cards from last week", "list cards in folder X"                           |


**Decision rules:**

- If the user provides 1–3 specific words that are clearly keywords → use `card list --keyword`.
- If the user describes an intent, asks a question, or uses a phrase that would benefit from semantic understanding → use `card rag --query`.
- If in doubt, prefer `card rag` — it is a superset that handles both precise and fuzzy queries well.
- You can combine: first `card rag` for discovery, then `card list --keyword` or `card list` with filters to narrow down.

## Step 2 — Refine the query before sending

Before calling `card rag --query`, **you must refine the user's raw input** into a better query:

1. **Extract core intent** — strip conversational filler ("hey can you", "I think I saved", "帮我找找").
2. **Expand with context** — if the user's wording is too short or vague, enrich with synonyms or related terms that better capture the intent. For example:
  - User: "database image upload" → Query: "Java implement database image upload and display on frontend"
  - User: "那篇关于网红的文章" → Query: "网红 流量 社会责任 互联网"
3. **Preserve language** — keep the query in the same language the user used. Do not translate unless asked.
4. **Keep it concise** — the refined query should be a single focused sentence or phrase, not a paragraph.

## Step 3 — Present results (stop here if sufficient)

After receiving RAG results:

1. **Summarize the list** to the user — show title, domain/source, and a one-line description for each result.
2. **If the user only needs to browse or locate bookmarks, stop here.** Do not fetch details unless needed.
3. Common "stop here" scenarios:
  - "帮我找找关于 X 的收藏" → show the list, done.
  - "I want to find that article about Y" → show the list, user picks one.
  - "有没有收藏过 Z 相关的东西" → show the list as confirmation.

## Step 4 — Fetch details progressively (only when needed)

Proceed to fetch `card detail` **only** when the user needs the actual content, for example:

- "帮我总结这些文章的要点"
- "Read the most relevant one for me"
- "把相关内容整理成一份报告"

**Progressive fetching rules:**

1. Do NOT fetch all details at once. Start with the most relevant 1–3 cards.
2. Use the RAG result order as initial relevance ranking.
3. Fetch one card's detail, evaluate if it answers the user's need. If yes, stop. If not, fetch the next.

## Step 5 — Re-rank by content relevance (only after detail fetch)

If you fetched details for multiple cards in Step 4:

1. **Compare each card's full content** against the user's original intent.
2. **Re-rank** by how well the content actually addresses the query — RAG ranking is based on metadata/embeddings and may not perfectly reflect full-content relevance.
3. **Filter out** cards whose content turned out to be irrelevant despite matching metadata.
4. **Present the re-ranked results** with brief explanations of why each is relevant.

## Example: Full RAG workflow

User: "帮我找找之前收藏的关于 Java 上传图片到数据库的教程，我需要参考一下具体实现"

Agent thinking:

- Intent: find tutorials about Java database image upload — needs semantic search → `card rag`
- The user says "参考具体实现" → will likely need detail content
- Refine query: "Java实现数据库图片上传功能 教程 前端渲染"

```bash
# Step 1-2: RAG search with refined query
cubox-cli card rag --query "Java实现数据库图片上传功能 教程 前端渲染"

# Step 3: Present results to user, e.g. "Found 5 related bookmarks: ..."
# User confirms: "第二篇看起来最相关，帮我看看"

# Step 4: Fetch only the selected card
cubox-cli card detail --id 7231297087109858478

# Step 5: (only if multiple details fetched) Re-rank and present
```

