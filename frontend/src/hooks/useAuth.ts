// hooks/useAuth.ts

import { useEffect, useState } from 'react';
import { useRouter } from 'next/router';

const useAuth = (user: any) => {
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  useEffect(() => {
    // If the user is not authenticated, redirect to the login page
    if (!user) {
      router.push('/login');
    }
    setLoading(false);
  }, [user, router]);

  return loading;
};

export default useAuth;
