'use client';
import { useUser } from '@clerk/nextjs';
import { Specta } from 'components/home/HomePage';
import LandingPage from 'components/home/LandingPage';


const HomePage = () => {
  const { isLoaded, isSignedIn, user } = useUser();

  if (!isLoaded) {
    return <div>Loading...</div>;
  }

  return (
    <Specta></Specta>
    // <LandingPage
    //   user={isSignedIn ? { name: user?.fullName || 'User' } : undefined}
    // />
  );
};

export default HomePage;
