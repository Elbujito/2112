import { ReactNode } from 'react';

type IHeroOneButtonProps = {
  title: string;
  description?: string;
  button?: ReactNode;
};

const HeroOneButton = (props: IHeroOneButtonProps) => (
  <header className="text-center">
    <h1 className="text-5xl text-white-100 font-bold whitespace-pre-line leading-hero">
      {props.title}
    </h1>
    <div className="text-4xl mt-4 mb-16">{props.description}</div>
    {props.button}
  </header>
);

export { HeroOneButton };
