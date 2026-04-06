import type { Provider } from "../hooks/useApiKey";

export interface TastingQuestion {
  question: string;
  choices: string[];
}

export interface TastingAnswers {
  [question: string]: string;
}

export async function generateTastingQuestions(
  provider: Provider,
  apiKey: string,
  model: string,
  freeText: string,
): Promise<TastingQuestion[]> {
  const prompt = `あなたはコーヒーの専門家です。以下のコーヒーの感想を読んで、テイスティングノートをより詳しくするための深掘り質問を3〜5問、選択肢付きで生成してください。

感想: "${freeText}"

必ずJSON配列で返してください。フォーマット:
[
  {
    "question": "質問文",
    "choices": ["選択肢1", "選択肢2", "選択肢3"]
  }
]

JSON以外は一切出力しないでください。`;

  const response = await callAI(provider, apiKey, model, prompt);
  const jsonMatch = response.match(/\[[\s\S]*\]/);
  if (!jsonMatch) throw new Error("質問の生成に失敗しました");
  return JSON.parse(jsonMatch[0]) as TastingQuestion[];
}

export async function generateTastingNote(
  provider: Provider,
  apiKey: string,
  model: string,
  freeText: string,
  answers: TastingAnswers,
): Promise<string> {
  const answerText = Object.entries(answers)
    .map(([q, a]) => `- ${q}: ${a}`)
    .join("\n");

  const prompt = `あなたはコーヒーの専門家です。以下の情報をもとに、プロフェッショナルなテイスティングノートを200字以内で書いてください。

最初の感想: "${freeText}"

深掘り回答:
${answerText}

【重要】回答がない項目については推測・補完せず、提供された情報のみに基づいて記述してください。
テイスティングノートのみを出力してください（見出しや説明は不要）:`;

  return callAI(provider, apiKey, model, prompt);
}

export async function generateRecommendComment(
  provider: Provider,
  apiKey: string,
  model: string,
  score: number,
  origin: string,
  name: string,
): Promise<string> {
  const prompt = `あなたはコーヒーの専門家です。以下のコーヒーについて、ユーザーの好みスコア（${score.toFixed(1)}/5.0）をもとに、一言コメントを50字以内で書いてください。

産地: ${origin || "不明"}
豆の名前: ${name || "不明"}

コメントのみを出力:`;

  return callAI(provider, apiKey, model, prompt);
}

export async function generateProfileComment(
  provider: Provider,
  apiKey: string,
  model: string,
  topOrigins: string[],
  favoriteRoast: string,
  flavorWords: string[],
): Promise<string> {
  const prompt = `あなたはコーヒーの専門家です。以下のユーザーの嗜好データをもとに、このユーザーのコーヒーの好みを表す一言コメントを80字以内で書いてください。

好きな産地TOP3: ${topOrigins.join("、") || "データなし"}
よく飲む焙煎度: ${favoriteRoast || "データなし"}
好きなフレーバーワード: ${flavorWords.slice(0, 5).join("、") || "データなし"}

コメントのみを出力:`;

  return callAI(provider, apiKey, model, prompt);
}

async function callAI(
  provider: Provider,
  apiKey: string,
  model: string,
  prompt: string,
): Promise<string> {
  switch (provider) {
    case "claude":
      return callClaude(apiKey, model, prompt);
    case "openai":
      return callOpenAI(apiKey, model, prompt);
    case "google":
      return callGoogle(apiKey, model, prompt);
  }
}

async function callClaude(apiKey: string, model: string, prompt: string): Promise<string> {
  const res = await fetch("https://api.anthropic.com/v1/messages", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "x-api-key": apiKey,
      "anthropic-version": "2023-06-01",
      "anthropic-dangerous-direct-browser-access": "true",
    },
    body: JSON.stringify({
      model,
      max_tokens: 1024,
      messages: [{ role: "user", content: prompt }],
    }),
  });

  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(
      (err as { error?: { message?: string } })?.error?.message ?? `API error ${res.status}`,
    );
  }

  const data = (await res.json()) as {
    content: Array<{ type: string; text: string }>;
  };
  return (data.content.find((c) => c.type === "text")?.text ?? "").trim();
}

async function callOpenAI(apiKey: string, model: string, prompt: string): Promise<string> {
  const res = await fetch("https://api.openai.com/v1/chat/completions", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${apiKey}`,
    },
    body: JSON.stringify({
      model,
      max_tokens: 1024,
      messages: [{ role: "user", content: prompt }],
    }),
  });

  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(
      (err as { error?: { message?: string } })?.error?.message ?? `API error ${res.status}`,
    );
  }

  const data = (await res.json()) as {
    choices: Array<{ message: { content: string } }>;
  };
  return (data.choices[0]?.message.content ?? "").trim();
}

async function callGoogle(apiKey: string, model: string, prompt: string): Promise<string> {
  const url = `https://generativelanguage.googleapis.com/v1beta/models/${model}:generateContent?key=${apiKey}`;
  const res = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      contents: [{ parts: [{ text: prompt }] }],
    }),
  });

  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(
      (err as { error?: { message?: string } })?.error?.message ?? `API error ${res.status}`,
    );
  }

  const data = (await res.json()) as {
    candidates: Array<{ content: { parts: Array<{ text: string }> } }>;
  };
  return (data.candidates[0]?.content.parts[0]?.text ?? "").trim();
}
