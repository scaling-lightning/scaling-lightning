// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require('prism-react-renderer/themes/github');
const darkCodeTheme = require('prism-react-renderer/themes/dracula');

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'Scaling Lightning',
  tagline: 'A Testing Toolkit for the Lightning Network',
  favicon: 'img/favicon.ico',

  // Set the production url of your site here
  url: 'https://scalinglightning.com',
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: '/',

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: 'scaling-lightning', // Usually your GitHub org/user name.
  projectName: 'scaling-lightning', // Usually your repo name.

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  // Even if you don't use internalization, you can use this field to set useful
  // metadata like html lang. For example, if your site is Chinese, you may want
  // to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        // docs: {
        //   sidebarPath: require.resolve('./sidebars.js'),
        //   // Please change this to your repo.
        //   // Remove this to remove the "edit this page" links.
        //   editUrl:
        //     'https://github.com/facebook/docusaurus/tree/main/packages/create-docusaurus/templates/shared/',
        // },
        blog: {
          showReadingTime: true,
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      colorMode: {
        defaultMode: 'light',
        disableSwitch: true,
        respectPrefersColorScheme: false,
      },
      // Replace with your project's social card
      image: 'img/scaling-lightning-social-card.png',
      navbar: {
        title: 'Scaling Lightning',
        logo: {
          alt: 'Scaling Lightning Logo',
          src: 'img/scaling-lightning-circle.png',
        },
        items: [
          // {
          //   type: 'docSidebar',
          //   sidebarId: 'tutorialSidebar',
          //   position: 'left',
          //   label: 'Tutorial',
          // },
          { to: '/blog', label: 'Blog', position: 'left' },
          {
            href: 'https://github.com/scaling-lightning/scaling-lightning',
            label: 'GitHub',
            position: 'right',
          },
        ],
      },
      footer: {
        style: 'light',
        links: [
          {
            title: 'Resources',
            items: [
              {
                label: 'GitHub',
                to: 'https://github.com/scaling-lightning/scaling-lightning',
              },
            ],
          },
          {
            title: 'Community',
            items: [
              {
                to: 'https://twitter.com/scalingln',
                label: 'Twitter',
              },
              {
                to: 'https://t.me/+AytRsS0QKH5mMzM8',
                label: 'Telegram',
              },
            ],
          },
          {
            title: `Made with 🧡 by Bitcoiners`,
          },
        ],
        // copyright: `Made with 🧡 by Bitcoiners`,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
      },
    }),
};

module.exports = config;
