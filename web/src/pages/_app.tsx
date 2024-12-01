import { AppProps } from 'next/app';
import 'animate.css/animate.min.css';
import '../styles/global.css';

const MyApp = ({ Component, pageProps }: AppProps) => (
  <Component {...pageProps} />
);

export default MyApp;
