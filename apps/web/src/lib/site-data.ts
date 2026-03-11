export type Auction = {
  id: string;
  title: string;
  office: string;
  category: string;
  closingAt: string;
  viewingAt: string;
  priceBand: string;
  warnings: string[];
  officialUrl: string;
  summary: string;
};

export type Article = {
  slug: string;
  title: string;
  summary: string;
  sections: { heading: string; body: string }[];
};

export type CommunityPost = {
  id: string;
  title: string;
  office: string;
  author: string;
  body: string;
};

export type Advisor = {
  id: string;
  name: string;
  specialty: string;
  note: string;
  responseTime: string;
};

export const disclaimers = [
  "現狀交付，不負瑕疵擔保",
  "官方標售不適用七日鑑賞期",
  "得標價通常不含營業稅與其他規費",
];

export const auctions: Auction[] = [
  {
    id: "lot-camera-001",
    title: "沒入數位相機與鏡頭一批",
    office: "臺北關",
    category: "3C 電子",
    closingAt: "2026-03-16 14:00",
    viewingAt: "2026-03-15 10:00",
    priceBand: "NT$50,000 - 80,000",
    warnings: ["現狀交付", "不負瑕疵擔保"],
    officialUrl: "https://web.customs.gov.tw/taipei/htmlList/86208152b05545bdad39749ea730870d",
    summary: "熱門相機批次，適合二手設備商與攝影工作室關注。",
  },
  {
    id: "lot-car-002",
    title: "進口車零組件與檢附文件一批",
    office: "臺中關",
    category: "交通工具",
    closingAt: "2026-03-20 15:30",
    viewingAt: "2026-03-19 11:00",
    priceBand: "NT$150,000 - 280,000",
    warnings: ["需確認驗車與牌照流程"],
    officialUrl: "https://web.customs.gov.tw/taichung",
    summary: "含進口車相關零件，適合有驗車與文件處理經驗者。",
  },
  {
    id: "lot-liquor-003",
    title: "菸酒混合拍賣品",
    office: "高雄關",
    category: "特殊管制",
    closingAt: "2026-03-22 09:30",
    viewingAt: "2026-03-21 13:30",
    priceBand: "NT$90,000 - 140,000",
    warnings: ["需確認菸酒相關資格"],
    officialUrl: "https://web.customs.gov.tw/kaohsiung/htmlList/127d64b348734567b41bdc83590b9d12",
    summary: "屬管制品項，需在投標前確認資格與後續處理流程。",
  },
];

export const articles: Article[] = [
  {
    slug: "bid-form-guide",
    title: "新手指南：標單怎麼填",
    summary: "拆解標單欄位、信封格式與常見退件原因。",
    sections: [
      {
        heading: "標單內容",
        body: "確認標號、投標金額、大寫金額與聯絡資料一致，避免空白欄位與塗改未蓋章。",
      },
      {
        heading: "押標金流程",
        body: "先確認公告金額與繳納方式，保留匯款憑證，並在截止前完成通信或現場投遞。",
      },
    ],
  },
  {
    slug: "deposit-and-tax",
    title: "押標金與稅費試算",
    summary: "用得標價、營業稅與額外規費快速估算總成本。",
    sections: [
      {
        heading: "營業稅",
        body: "先以得標價乘以 5% 作為基礎估算，再視公告備註加入貨物稅或服務費。",
      },
      {
        heading: "提領成本",
        body: "別忽略倉儲、搬運、驗車與補件時間成本，特殊管制品項尤其要預留額外預算。",
      },
    ],
  },
];

export const posts: CommunityPost[] = [
  {
    id: "post-001",
    title: "臺北關相機批次看貨紀錄",
    office: "臺北關",
    author: "攝影器材商",
    body: "多數鏡頭有使用痕跡，外盒缺件比率高，建議把看貨時間拉滿。",
  },
  {
    id: "post-002",
    title: "高雄關菸酒案提醒",
    office: "高雄關",
    author: "物流從業者",
    body: "公告寫得簡短，但實際上資格確認是重點，沒有相關證照不要貿然出價。",
  },
];

export const advisors: Advisor[] = [
  {
    id: "advisor-001",
    name: "王顧問",
    specialty: "進口車標售流程",
    note: "熟悉驗車、牌照與報關補件流程。",
    responseTime: "平均 4 小時回覆",
  },
  {
    id: "advisor-002",
    name: "林代標",
    specialty: "高價 3C 與精品批次",
    note: "偏重看貨估值與風險評估。",
    responseTime: "平均 1 個工作天回覆",
  },
];

export const courses = [
  {
    id: "course-001",
    title: "如何從標售貨物中獲利",
    summary: "從看貨清單、風險折價到轉售節奏的完整框架。",
  },
  {
    id: "course-002",
    title: "進口車標售實務",
    summary: "檢附文件、驗車、領牌與轉售注意事項。",
  },
];

export const historicalCoverageStart = "2026-03-11";
