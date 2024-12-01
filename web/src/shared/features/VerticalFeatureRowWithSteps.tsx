import className from 'classnames';
import { ReactNode } from 'react';

type IVerticalRuleRowProps = {
  title: string;
  description: string;
  reverse?: boolean;
  steps?: ReactNode;
};

const VerticalFeatureRowWithSteps = (props: IVerticalRuleRowProps) => {
  const verticalFeatureClass = className(
    'mt-20',
    'flex',
    'flex-wrap',
    'items-center',
    {
      'flex-row-reverse': props.reverse,
    }
  );

  return (
    <div className={verticalFeatureClass}>
      <div className=" w-full text-left sm:px-6">
        <h3 className="text-3xl text-gray-900 font-semibold">{props.title}</h3>
        <div className="mt-6 text-gray-900 text-xl leading-9">{props.description}</div>
      </div>
      <div className="text-gray-900 text-xl w-full text-left sm:px-6 leading-9">{props.steps}</div>
    </div>
  );
};

export { VerticalFeatureRowWithSteps };
