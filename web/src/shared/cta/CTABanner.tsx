import { ReactNode } from 'react';


type ICTABannerProps = {
  title?: string;
  subtitle?: ReactNode;
  button?: ReactNode;
};

const CTABanner = (props: ICTABannerProps) => (
  <div>
  <div className="text-center flex flex-col p-4 sm:text-left sm:items-center sm:justify-between sm:p-12 bg-white-100 rounded-md">
    <div className="text-center flex flex-row p-4 sm:text-left sm:flex-row">
    <div className="text-2xl font-semibold">
      <div className="text-black-100">{props.title}</div>
      <div className="text-primary-800">{props.subtitle}</div>
    </div>
    </div>
  </div>
  </div>
);

export { CTABanner };
