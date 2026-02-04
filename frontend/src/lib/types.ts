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
    transparent: boolean;
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
  total_coin?: number;
  price?: number;
  coin_type?: string;
  total_num?: number;
  combo_num?: number;
  batch_combo_id?: string;
  batch_combo_num?: number;
  dmscore?: number;
  gift_id?: number;
  is_join_receiver?: boolean;
  is_naming?: boolean;
  is_show?: number;
  name_color?: string;
  r_uname?: string;
  receive_user_info?: {
    uid: number;
    uname: string;
  };
  ruid?: number;
  // 内部使用的字段
  _rawNum?: number;
  _rawComboNum?: number;
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
