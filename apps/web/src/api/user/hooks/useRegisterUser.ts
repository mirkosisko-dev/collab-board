import { useMutation } from '@tanstack/react-query';
import userApi from '../service';
import type { RegisterPayload } from '../types';

const useRegisterUser = () =>
	useMutation({
		mutationFn: (payload: RegisterPayload) => userApi.register(payload),
	});

export default useRegisterUser;
