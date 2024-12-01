import { Meta } from '../../shared/layout/Meta';
import { AppConfig } from '../../config/AppConfig';
import { Footer } from '../footer/Footer';
import { Hero } from '../hero/Hero';
import { Rules } from './Rules';

const Rulebook = () => (
  <div className="antialiased text-white-100">
    <Meta title={AppConfig.title} description={AppConfig.description} />
    <Hero title='SpaceOps Rulebook'/>
    <Rules />
    <Footer />
  </div>
);

export { Rulebook };
