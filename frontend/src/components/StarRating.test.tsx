import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { StarRating } from "./StarRating";

describe("StarRating", () => {
  it("5つの星を描画する", () => {
    render(<StarRating value={3} />);
    const stars = screen.getAllByText("★");
    expect(stars).toHaveLength(5);
  });

  it("value以下の星を金色で表示する", () => {
    render(<StarRating value={3} />);
    const stars = screen.getAllByText("★");
    expect(stars[0]).toHaveStyle({ color: "#c47a1b" });
    expect(stars[1]).toHaveStyle({ color: "#c47a1b" });
    expect(stars[2]).toHaveStyle({ color: "#c47a1b" });
    expect(stars[3]).toHaveStyle({ color: "#d4c5a9" });
    expect(stars[4]).toHaveStyle({ color: "#d4c5a9" });
  });

  it("クリックで onChange が呼ばれる", async () => {
    const onChange = vi.fn();
    render(<StarRating value={0} onChange={onChange} />);
    const stars = screen.getAllByText("★");
    await userEvent.click(stars[2]);
    expect(onChange).toHaveBeenCalledWith(3);
  });

  it("readonly の場合 onChange が呼ばれない", async () => {
    const onChange = vi.fn();
    render(<StarRating value={0} onChange={onChange} readonly />);
    const stars = screen.getAllByText("★");
    await userEvent.click(stars[2]);
    expect(onChange).not.toHaveBeenCalled();
  });

  it("size に応じてフォントサイズが変わる", () => {
    const { rerender } = render(<StarRating value={1} size="sm" />);
    expect(screen.getAllByText("★")[0]).toHaveStyle({ fontSize: "1.1rem" });

    rerender(<StarRating value={1} size="lg" />);
    expect(screen.getAllByText("★")[0]).toHaveStyle({ fontSize: "2rem" });
  });
});
