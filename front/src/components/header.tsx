import { useState } from "react";
import { getGlobalProps } from "../lib/props";
import BracketButton from "./bracket_button";
import "./header.css";

export default function Header() {
  const { IsLoggedIn } = getGlobalProps();

  const navigationForLoggedInUser = [
    {
      label: "一覧",
      to: "/all",
    },
    {
      label: "新規作成",
      to: "/new",
    },
  ];
  const navigationForAnonymousUser = [
    {
      label: "一覧",
      to: "/all",
    },
  ];

  const loginButton = {
    label: "ログイン",
    to: "/login",
  };
  const logoutButton = {
    label: "ログアウト",
    to: "/logout",
  };

  const [isOpenSubHeader, setIsOpenSubHeader] = useState(false);

  return (
    <div className="components-header">
      <div className="primary">
        <Breadcrumbs
          primary
          onClick={() => setIsOpenSubHeader(!isOpenSubHeader)}
        />
        <div className="nav wide">
          {IsLoggedIn ? (
            <BracketButton values={navigationForLoggedInUser} />
          ) : (
            <BracketButton values={navigationForAnonymousUser} />
          )}
          {IsLoggedIn ? (
            <BracketButton values={[logoutButton]} />
          ) : (
            <BracketButton values={[loginButton]} />
          )}
        </div>
        <div className="nav narrow">
          {IsLoggedIn ? (
            <BracketButton values={[{ label: "新規作成", to: "/new" }]} />
          ) : (
            <BracketButton values={[loginButton]} />
          )}
        </div>
      </div>
      <div
        className={
          "wrapped wrapped-breadcrumbs " + (isOpenSubHeader ? "show" : "")
        }
      >
        <Breadcrumbs primary={false} />
      </div>
      {IsLoggedIn && (
        <div
          className={"wrapped wrapped-nav " + (isOpenSubHeader ? "show" : "")}
        >
          <BracketButton values={[{ label: "一覧", to: "/all" }]} />
          <BracketButton values={[logoutButton]} />
        </div>
      )}
    </div>
  );
}

function Breadcrumbs({
  primary,
  onClick,
}: {
  primary: boolean;
  onClick?: () => void;
}) {
  const { Breadcrumbs } = getGlobalProps();
  return (
    <div
      onClick={onClick}
      className={`${primary ? "primary-breadcrumbs" : "secondary"} breadcrumbs`}
    >
      <ul>
        {Breadcrumbs.map((v) => (
          <li key={v.Location}>
            <a href={v.Location}>{v.Label}</a>
          </li>
        ))}
      </ul>
    </div>
  );
}
