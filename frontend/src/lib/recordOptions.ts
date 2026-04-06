export type RoastLevel = "light" | "medium" | "dark";
export type BrewMethod = "drip" | "espresso" | "french_press" | "other";

export const ROAST_OPTIONS: Array<{ value: RoastLevel; label: string }> = [
  { value: "light", label: "浅煎り" },
  { value: "medium", label: "中煎り" },
  { value: "dark", label: "深煎り" },
];

export const BREW_OPTIONS: Array<{ value: BrewMethod; label: string }> = [
  { value: "drip", label: "ドリップ" },
  { value: "espresso", label: "エスプレッソ" },
  { value: "french_press", label: "フレンチプレス" },
  { value: "other", label: "その他" },
];

export const ROAST_LABEL: Record<string, string> = {
  light: "浅煎り",
  medium: "中煎り",
  dark: "深煎り",
};

export const BREW_LABEL: Record<string, string> = {
  drip: "ドリップ",
  espresso: "エスプレッソ",
  french_press: "フレンチプレス",
  other: "その他",
};

export interface RecordFormValues {
  name: string;
  rating: number;
  origin: string;
  roastLevel: RoastLevel | "";
  shop: string;
  price: string;
  purchasedAt: string;
  tastingNote: string;
  brewMethod: BrewMethod | "";
}
