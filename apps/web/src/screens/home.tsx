import Login from '@/components/login';
import Register from '@/components/register';
import React from 'react';

const Home = () => {
	const [isLogin, setIsLogin] = React.useState(true);

	return isLogin ? (
		<Login switchToRegister={() => setIsLogin(false)} />
	) : (
		<Register switchToLogin={() => setIsLogin(true)} />
	);
};

export default Home;
