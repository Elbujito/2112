'use client';

import Link from 'next/link';
import { useTheme } from 'next-themes';
import ThemeSwitch from 'components/shared/ThemeSwitch';
import { useState, useEffect } from 'react';

const LandingHeader = ({
  className,
  logo,
  logoDark,
  hideMenuItems,
}: {
  className?: string;
  logo: React.ReactNode;
  logoDark: React.ReactNode;
  hideMenuItems?: boolean;
}) => {
  const { theme = 'system', setTheme } = useTheme();
  const [mounted, setMounted] = useState(false);

  const basicNavLinks = [
    { title: 'Home', href: '/' },
    { title: 'Features', href: '/features' },
    { title: 'Pricing', href: '/pricing' },
    { title: 'About', href: '/about' },
    { title: 'Contact', href: '/contact' },
  ];


  useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted) {
    return null;
  }

  return (
    <header
      className={`flex items-center justify-between py-10 flex-wrap w-full mb-20 lg:mb-32 pt-6 wide-container ${className}`}
    >
      <div>
        <Link href="/" aria-label="Home">
          {theme === 'dark' ? <div>{logoDark}</div> : <div>{logo}</div>}
        </Link>
      </div>
      <div className="flex items-center leading-5 gap-4 sm:gap-6">
        {!hideMenuItems &&
          basicNavLinks.map((link) => (
            <Link
              key={link.title}
              href={link.href}
              className="nav-link hidden sm:block hover:text-gray-700 dark:hover:text-gray-300"
            >
              {link.title}
            </Link>
          ))}
        <ThemeSwitch />
        {!hideMenuItems && (
          <button
            aria-label="Toggle mobile menu"
            className="block sm:hidden p-2"
          >
            <span>☰</span>
          </button>
        )}
      </div>
    </header>
  );
};

export default LandingHeader;
