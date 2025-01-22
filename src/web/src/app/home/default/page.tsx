'use client';
import { useUser } from '@clerk/nextjs';
import { Specta } from 'components/home/HomePage';


const HomePage = () => {
  const { isLoaded, isSignedIn, user } = useUser();

  if (!isLoaded) {
    return <div>Loading...</div>;
  }

  return (
    <Specta></Specta>
  );
};

export default HomePage;
