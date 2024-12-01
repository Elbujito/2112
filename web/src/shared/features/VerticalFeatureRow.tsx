import className from 'classnames';
import { useRouter } from 'next/router';
import Image from'next/image';
import { AnimationOnScroll } from 'react-animation-on-scroll';

type IVerticalFeatureRowProps = {
  title: string;
  description: string;
  image?: string;
  imageAlt?: string;
  reverse?: boolean;
};


export default function VerticalFeatureRow(props: IVerticalFeatureRowProps) {

  const verticalFeatureClass = className(
    'mt-20',
    'flex',
    'flex-wrap',
    'items-center',
    {
      'flex-row-reverse': props.reverse,
    }
  );

  const router = useRouter();
  
  return (

    <div className={verticalFeatureClass}>
            <div className=" w-full sm:w-1/2 text-center sm:px-6">
            { props.reverse ? <AnimationOnScroll animateOnce={true} animateIn="animate__fadeInRightBig">

        <h3 className="text-3xl text-gray-900 font-semibold">{props.title}</h3>
        <div className="mt-6 text-gray-900 text-xl leading-9">{props.description}</div>
        </AnimationOnScroll> : 
        <AnimationOnScroll animateOnce={true} animateIn="animate__fadeInLeftBig">

        <h3 className="text-3xl text-gray-900 font-semibold">{props.title}</h3>
        <div className="mt-6 text-gray-900 text-xl leading-9">{props.description}</div>
        </AnimationOnScroll>}
      </div>


      <div className="w-full sm:w-1/2 p-6">
      { props.reverse ?   <AnimationOnScroll animateOnce={true} animateIn="animate__zoomIn">
        <Image width={500} height={500} src={`${router.basePath}${props.image}`} alt={props.imageAlt} />
        </AnimationOnScroll> : <AnimationOnScroll animateOnce={true} animateIn="animate__zoomIn">
        <Image width={500} height={500} src={`${router.basePath}${props.image}`} alt={props.imageAlt} />
        </AnimationOnScroll> }
      </div>

    </div>

  );
};
