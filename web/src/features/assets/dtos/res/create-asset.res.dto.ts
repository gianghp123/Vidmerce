export type UploadInfo = {
	fileKey:   string;
	uploadUrl: string;
	expiresIn: number;
};

export type CreateAssetRes = {
	assetId: string;
  uploads: UploadInfo[];
};
