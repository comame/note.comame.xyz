import { useState } from "react";
import BracketButton from "./bracket_button";
import "./header.css";
import { useMediaQuery } from "../lib/useMediaQuery";

type props = {
  isLoggedIn: boolean;
  breadcrumbs: {
    label: string;
    location: string;
  }[];
};

export default function Header({ isLoggedIn, breadcrumbs }: props) {
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
          breadcrumbs={breadcrumbs}
        />
        <div className="nav wide">
          {isLoggedIn ? (
            <BracketButton values={navigationForLoggedInUser} />
          ) : (
            <BracketButton values={navigationForAnonymousUser} />
          )}
          {isLoggedIn ? (
            <BracketButton values={[logoutButton]} />
          ) : (
            <BracketButton values={[loginButton]} />
          )}
        </div>
        <div className="nav narrow">
          {isLoggedIn ? (
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
        <Breadcrumbs primary={false} breadcrumbs={breadcrumbs} />
      </div>
      {isLoggedIn && (
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
  breadcrumbs,
}: {
  primary: boolean;
  onClick?: () => void;
  breadcrumbs: {
    label: string;
    location: string;
  }[];
}) {
  const isNarrowWidth = useMediaQuery("(max-width: 550px)");

  return (
    <div
      onClick={onClick}
      className={`${primary ? "primary-breadcrumbs" : "secondary"} breadcrumbs`}
    >
      <ul>
        {breadcrumbs.map((v) => (
          <li key={v.location}>
            <a href={v.location} inert={isNarrowWidth && primary}>
              {v.label}
            </a>
          </li>
        ))}
      </ul>
    </div>
  );
}
