import Link from 'next/link';

import { Background } from '../../shared/background/Background';
import { HeroOneButton } from '../../shared/hero/HeroOneButton';
import { Section } from '../../shared/layout/Section';
import { NavbarTwoColumns } from '../../shared/navigation/NavbarTwoColumns';
import { Logo } from '../../shared/logo/Logo';

type IHeroProps = {
  title: string;
  description?: string;
};

const Hero = (props: IHeroProps) => {

  return (
  <Background color="bg-black-100" >
     <Section styles={` py-2 px-3 `}> 
      <NavbarTwoColumns color="text-gray-800" childrenRight={
        <NavbarTwoColumns color="text-white-100" logo={<span></span>}>
        <li>
          <Link  href="/">
            <a>Home</a>
          </Link>
        </li>
        <li>
          <Link href="/contact">
            <a>Contact</a>
          </Link>
        </li>
      </NavbarTwoColumns>
       }  logo={<Logo textColor="text-white-100" logoColor="text-white-100" xl />}>
      </NavbarTwoColumns>
    </Section> 

    <Section styles={`pt-2 pb-2 max-w-screen-lg mx-auto px-3 `}>
      <HeroOneButton
        title={props.title}
        description={props.description}
      />
    </Section>
  </Background>
    );
};

export { Hero };
