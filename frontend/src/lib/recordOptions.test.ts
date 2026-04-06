import { describe, it, expect } from "vitest";
import { ROAST_OPTIONS, BREW_OPTIONS, ROAST_LABEL, BREW_LABEL } from "./recordOptions";

describe("ROAST_OPTIONS", () => {
  it("3種類の焙煎レベルを含む", () => {
    expect(ROAST_OPTIONS).toHaveLength(3);
  });

  it("値とラベルが正しい", () => {
    expect(ROAST_OPTIONS[0]).toEqual({ value: "light", label: "浅煎り" });
    expect(ROAST_OPTIONS[1]).toEqual({ value: "medium", label: "中煎り" });
    expect(ROAST_OPTIONS[2]).toEqual({ value: "dark", label: "深煎り" });
  });
});

describe("BREW_OPTIONS", () => {
  it("4種類の抽出方法を含む", () => {
    expect(BREW_OPTIONS).toHaveLength(4);
  });

  it("値とラベルが正しい", () => {
    expect(BREW_OPTIONS[0]).toEqual({ value: "drip", label: "ドリップ" });
    expect(BREW_OPTIONS[1]).toEqual({ value: "espresso", label: "エスプレッソ" });
    expect(BREW_OPTIONS[2]).toEqual({ value: "french_press", label: "フレンチプレス" });
    expect(BREW_OPTIONS[3]).toEqual({ value: "other", label: "その他" });
  });
});

describe("ROAST_LABEL", () => {
  it("各焙煎レベルを日本語に変換する", () => {
    expect(ROAST_LABEL["light"]).toBe("浅煎り");
    expect(ROAST_LABEL["medium"]).toBe("中煎り");
    expect(ROAST_LABEL["dark"]).toBe("深煎り");
  });
});

describe("BREW_LABEL", () => {
  it("各抽出方法を日本語に変換する", () => {
    expect(BREW_LABEL["drip"]).toBe("ドリップ");
    expect(BREW_LABEL["espresso"]).toBe("エスプレッソ");
    expect(BREW_LABEL["french_press"]).toBe("フレンチプレス");
    expect(BREW_LABEL["other"]).toBe("その他");
  });
});
