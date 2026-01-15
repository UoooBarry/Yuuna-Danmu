export interface ServerSettings {
    name: string;
    type: string;
    port: number;
    enabled: boolean;
}

export interface AppConfig {
    room_id: number;
    cookie: string;
    refresh_token: string;
    servers: ServerSettings[];
}

export interface DanmuMsg {
    type: 'danmu';
    medalName?: string;
    medalLevel?: number;
    nickname: string;
    content: string;
}

export interface GiftMsg {
    type: 'gift';
    combo_id?: string;
    gift_num: number;
    face?: string;
    medalName?: string;
    medalLevel?: number;
    uname?: string;
    action?: string;
    gift_name?: string;
    gift_info?: {
        gif: string;
    };
    combo_send?: {
        combo_id: string;
        combo_num: number;
    };
    combo_total_coin?: number;
    coin_type?: string;
}

export type DanmuItem = DanmuMsg | GiftMsg;

export interface SuperChatMsg {
    start_time: number;
    end_time: number;
    price: number;
    message: string;
    user_info: {
        face: string;
        uname: string;
    };
}
