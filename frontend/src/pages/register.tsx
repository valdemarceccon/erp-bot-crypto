import { NextPage } from 'next';
import RegisterForm from '../components/RegisterForm';
import MainLayout from '@/layouts/MainLayout';

const RegisterPage: NextPage = () => {
  const handleRegister = (email: string, name: string, password: string) => {
    console.log('Email:', email, 'Name:', name, 'Password:', password);
  };

  return (
    <MainLayout loggedIn={false}>
      <RegisterForm onRegister={handleRegister} />
    </MainLayout>
  );
};

export default RegisterPage;
