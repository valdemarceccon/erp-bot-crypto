import type { NextPage } from 'next';
import MainLayout from '../layouts/MainLayout';
import LoginForm from '../components/LoginForm';

const LoginPage: NextPage = () => {
  const handleLogin = (email: string, password: string) => {
    console.log('Email:', email);
    console.log('Password:', password);
  };

  return (
    <MainLayout loggedIn={false}>
      <LoginForm onLogin={handleLogin} />
    </MainLayout>
  );
};

export default LoginPage;

