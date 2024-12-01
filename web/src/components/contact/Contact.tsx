import { Meta } from '../../shared/layout/Meta';
import { AppConfig } from '../../config/AppConfig';
import { ContactForm } from './ContactForm';
import { Footer } from '../footer/Footer';
import { Hero } from '../hero/Hero';


const Contact = () => {
  return (
  <div className="antialiased text-white-100">
    <Meta title={AppConfig.title} description={AppConfig.description} />
    <Hero title="" description="Contact Information" />
    <div className="py-6 max-w-screen-lg mx-auto px-3 "
  >
    <ContactForm ></ContactForm>
    </div>
    <Footer />
  </div>
  );
};

export { Contact };
