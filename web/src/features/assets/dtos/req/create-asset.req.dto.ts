export type CreateAssetReqDto = {
	imageCount: number;
	name: string;
	price: number;
	productUrl: string;
};

export type ConfirmUploadDto = {
	name: string;
	price: number;
	productUrl: string;
	fileKeys: string[];
};
