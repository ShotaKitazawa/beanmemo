import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { RoastBar } from "./RoastBar";

describe("RoastBar", () => {
  it("データがない場合「データなし」を表示する", () => {
    render(<RoastBar byRoastLevel={[]} />);
    expect(screen.getByText("データなし")).toBeInTheDocument();
  });

  it("全件数が0の場合「データなし」を表示する", () => {
    render(
      <RoastBar
        byRoastLevel={[
          { label: "light", count: 0, avg_rating: 0 },
          { label: "medium", count: 0, avg_rating: 0 },
        ]}
      />,
    );
    expect(screen.getByText("データなし")).toBeInTheDocument();
  });

  it("各焙煎レベルの件数と平均評価を表示する", () => {
    render(
      <RoastBar
        byRoastLevel={[
          { label: "light", count: 2, avg_rating: 4.5 },
          { label: "medium", count: 3, avg_rating: 3.0 },
        ]}
      />,
    );
    expect(screen.getByText(/浅煎り: 2件 \(★4\.5\)/)).toBeInTheDocument();
    expect(screen.getByText(/中煎り: 3件 \(★3\.0\)/)).toBeInTheDocument();
  });

  it("件数が0の焙煎レベルはラベルを表示しない", () => {
    render(
      <RoastBar
        byRoastLevel={[
          { label: "light", count: 5, avg_rating: 4.0 },
          { label: "dark", count: 0, avg_rating: 0 },
        ]}
      />,
    );
    expect(screen.getByText(/浅煎り/)).toBeInTheDocument();
    expect(screen.queryByText(/深煎り/)).not.toBeInTheDocument();
  });

  it("バーの title にラベルと件数が表示される", () => {
    render(
      <RoastBar
        byRoastLevel={[
          { label: "light", count: 1, avg_rating: 5.0 },
          { label: "medium", count: 1, avg_rating: 4.0 },
        ]}
      />,
    );
    expect(screen.getByTitle("浅煎り: 1件")).toBeInTheDocument();
    expect(screen.getByTitle("中煎り: 1件")).toBeInTheDocument();
  });
});
