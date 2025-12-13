import type { LoginPayload, RegisterPayload } from './types';

const baseUrl = 'http://127.0.0.1:8080/api/v1';

const userApi = {
	register: async (payload: RegisterPayload) => {
		try {
			return await fetch(`${baseUrl}/register`, {
				method: 'POST',
				body: JSON.stringify(payload),
			});
		} catch (error) {
			throw new Error(error as string);
		}
	},

	login: async (payload: LoginPayload) => {
		try {
			await fetch(`${baseUrl}/login`, {
				method: 'POST',
				body: JSON.stringify(payload),
			});
		} catch (error) {
			throw new Error(error as string);
		}
	},
};

export default userApi;
