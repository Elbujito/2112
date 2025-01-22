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
    { title: 'Pricing', href: '/home/pricing' },
  ];

  useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted) {
    return null;
  }

  return (
    <header
      className={`flex items-center justify-between py-6 px-6 lg:px-12 w-full mb-20 lg:mb-32 wide-container ${className}`}
    >
      <div className="flex items-center gap-4">
        <Link href="/" aria-label="Home" className="flex items-center gap-2">
          {theme === 'dark' ? <div>{logoDark}</div> : <div>{logo}</div>}
          <span className="text-2xl font-bold tracking-wide text-gray-900 dark:text-gray-100">
            2112
          </span>
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
        {!hideMenuItems && (
          <Link
            href="/home/login"
            className="hidden sm:inline-block px-4 py-2 rounded-md text-white bg-blue-600 hover:bg-blue-500 dark:bg-blue-800 dark:hover:bg-blue-700 font-semibold shadow-md"
          >
            Sign In
          </Link>
        )}
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
