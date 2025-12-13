import { QueryClientProvider } from '@tanstack/react-query';
import { queryClient } from './config/query-client';
import Home from '@/screens/home';

function App() {
	return (
		<QueryClientProvider client={queryClient}>
			<div className="h-screen w-screen p-4">
				<Home />
			</div>
		</QueryClientProvider>
	);
}

export default App;
