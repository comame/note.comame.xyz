import "./bracket_button.css";

interface props {
  values: value[];
}

type value =
  | {
      label: string;
      to: string;
    }
  | {
      label: string;
      onClick: () => void;
    }
  | {
      label: string;
    };

export default function BracketButton({ values }: props) {
  return (
    <div className="components-bracket-button">
      {values.map((v) => (
        <span key={v.label} className="element">
          {"to" in v ? (
            <a href={v.to}>{v.label}</a>
          ) : "onClick" in v ? (
            <button type="button" onClick={v.onClick}>
              {v.label}
            </button>
          ) : (
            <span className="inert">{v.label}</span>
          )}
        </span>
      ))}
    </div>
  );
}
