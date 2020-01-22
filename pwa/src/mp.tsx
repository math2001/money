import React from "react";

interface Props {
  active: string;
  children: React.ReactElement<TabProps>[];
  onChange: (id: string) => void;
}

export class TabSet extends React.Component<Props> {
  render(): JSX.Element {
    return (
      <div>
        <nav className="tabset-nav">
          <ul>
            {this.props.children.map(tab => (
              <li key={tab.props.id}>
                <button
                  className={this.props.active === tab.props.id ? "active" : ""}
                  onClick={() => this.props.onChange(tab.props.id)}
                >
                  {tab.props.title}
                </button>
              </li>
            ))}
          </ul>
        </nav>
        {this.props.children.map(tab => {
          if (this.props.active === tab.props.id) {
            return tab;
          }
          return null;
        })}
      </div>
    );
  }
}

interface TabProps {
  id: string;
  title: string;
}

export class Tab extends React.Component<TabProps> {
  render() {
    return <section className="tabset-tab">{this.props.children}</section>;
  }
}
