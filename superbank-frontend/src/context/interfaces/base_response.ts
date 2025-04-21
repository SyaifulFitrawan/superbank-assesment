interface IMetaData {
    code: number;
    success: boolean;
    message: string;
}

interface IPaginator {
    itemCount: number;
    limit: number;
    pageCount: number;
    page: number;
    hasPrevPage: boolean;
    hasNextPage: boolean;
    prevPage: number | null;
    nextPage: number | null;
}

interface BaseResponse<T> {
    meta : IMetaData;
    data : T;
    paginator: IPaginator | undefined
}

export type {IMetaData, BaseResponse};
