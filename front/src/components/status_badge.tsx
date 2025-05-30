import "./status_badge.css";

type Props = {
  status: "public" | "private" | "url";
  inherit?: boolean;
  size: "s" | "m";
};

export default function StatusBadge({ status, inherit, size }: Props) {
  let label: string = status;
  if (size === "s") {
    switch (status) {
      case "public":
        label = "PU";
        break;
      case "private":
        label = "PR";
        break;
      case "url":
        label = "U";
        break;
    }
  }

  const className =
    `components-status-badge size-${size} status-${status} ` +
    (inherit ? "inherit " : " ");

  return <div className={className}>{label.toUpperCase()}</div>;
}
