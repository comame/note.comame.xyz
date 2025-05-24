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
    };

export default function BracketButton({ values }: props) {
  return (
    <div className="components-bracket-button">
      {values.map((v) => (
        <span key={v.label} className="element">
          {"to" in v ? (
            <a href={v.to}>{v.label}</a>
          ) : (
            <button onClick={v.onClick}>{v.label}</button>
          )}
        </span>
      ))}
    </div>
  );
}
