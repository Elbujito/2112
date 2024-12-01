import { AppConfig } from '../../config/AppConfig';

type ILogoProps = {
  xl?: boolean;
  textColor?: string;
  logoColor?: string;
};

const Logo = (props: ILogoProps) => {
  const size = props.xl ? '44' : '32';
  const fontStyle = props.xl
    ? 'font-semibold text-3xl'
    : 'font-semibold text-xl';

  return (
    <span className={`${props.textColor} inline-flex items-center ${fontStyle}`}>
      <img src="/assets/images/2112.png" alt="" width={size} height={size} />

      {AppConfig.site_name}
    </span>
  );
};

export { Logo };
