import { useEffect } from 'react';
import { useRouter } from 'next/router';
import { Typography } from '@mui/material';
import MainLayout from '@/layouts/MainLayout';
import useAuth from '@/hooks/useAuth';

const Home = () => {
  const loading = useAuth(null);

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <MainLayout loggedIn={false}>
      <Typography variant='h1'>Hello</Typography>
    </MainLayout>
  );
};

export default Home;
