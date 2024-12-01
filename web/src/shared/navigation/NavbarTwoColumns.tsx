import { ReactNode } from 'react';
import Link from 'next/link';

type INavbarProps = {
  logo?: ReactNode;
  children: ReactNode;
  childrenRight?: ReactNode;
  color?: string;
};

const NavbarTwoColumns = (props: INavbarProps) => (

  <div className={` ${props.logo} flex flex-wrap justify-between items-center`}>
    <div>
      <Link href="/">
        {props.logo}
      </Link>
    </div>
     

    <nav>
      <ul className={`navbar flex items-center font-medium text-xl   ${props.color} `}>
        {props.children}
      </ul>
    </nav>

    <nav>
      <ul className={`navbar flex items-center font-medium text-xl  ${props.color} `}>
        {props.childrenRight}
      </ul>
    </nav>

    <style jsx>
      {`
        .navbar :global(li:not(:first-child)) {
          @apply mt-0;
        }

        .navbar :global(li:not(:last-child)) {
          @apply mr-5;
        }
      `}
    </style>
  </div>
);

export { NavbarTwoColumns };
