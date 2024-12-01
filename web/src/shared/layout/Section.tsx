import { ReactNode } from 'react';
import { AnimationOnScroll } from 'react-animation-on-scroll';
type ISectionProps = {
  title?: string;
  description?: string;
  yPadding?: string;
  children: ReactNode;
  color?: string;
  styles?: string
};

const Section = (props: ISectionProps) => (
  <div
    className={`${props.styles}`}
  >
    {(props.title || props.description) && (
      <AnimationOnScroll animateOnce={true} animateIn="animate__fadeInDown">
      <div className="mb-12 text-center">
        {props.title && (
          <h2 className="text-4xl text-gray-900 font-bold">{props.title}</h2>
        )}
        {props.description && (
          <div className="  text-gray-900  mt-4 text-xl md:px-20">{props.description}</div>
        )}
      </div>
      </AnimationOnScroll>
    )}

    {props.children}
  </div>
);

export { Section };
