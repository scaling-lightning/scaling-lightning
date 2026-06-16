import React from 'react';
import clsx from 'clsx';
import styles from './styles.module.css';
import classNames from "classnames";

type FeatureItem = {
  title: string;
  Svg: React.ComponentType<React.ComponentProps<'svg'>>;
  description: JSX.Element;
};

const FeatureList = [
  {
    title: 'The Goal',
    Svg: require('@site/static/img/scaling-lightning-goal.svg').default,
    description: (
      <>
        This initiative aims to build a testing toolkit for the Lightning Network protocol, its implementations, and applications that depend on the Lightning Network.

        The goal is to collaborate as an industry to help scale the Lightning Network and the applications that depend on it.
      </>
    ),
  },
  {
    title: 'The Why',
    Svg: require('@site/static/img/scaling-lightning-strategy.svg').default,
    description: (
      <>
       Currently, there are unknowns and untested assumptions about how the Lightning Network and its applications will react to shocks in transaction volume, channels, nodes, gossip messages, etc. Having a set of tools and a signet Lightning Network will help Developers, Researchers, Operators and Novices.
      </>
    ),
  },
  {
    title: 'The How',
    Svg: require('@site/static/img/scaling-lightning-idea.svg').default,
    description: (
      <>
        We are still in the early stages of planning, but the first tool we are building will be a tool to quickly generate one or more Lightning Nodes. These nodes can connect either to a public signet Lightning Network or a private Regtest Lightning Network for any combination of LN implementations (CLN, LND, LDK, Acinq etc.).
      </>
    ),
  },
];

function Feature({title, Svg, description}: FeatureItem) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
        <Svg className={styles.featureSvg} role="img" />
      </div>
      <div className={classNames("text--center padding-horiz--md", styles.splashBanner)}>
        <h3>{title}</h3>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): JSX.Element {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
