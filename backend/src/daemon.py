import datetime
from typing import Dict
from typing import List

from pybit import position
from pybit.unified_trading import HTTP
from pydantic import BaseModel
from src.dependencies.database import SessionLocal
from src.models.user import ApiKey
from src.models.user import ApiKeyStatusEnum
from src.models.user import ClosedPnl
from src.models.user import User


class TradeOrder(BaseModel):
    symbol: str
    orderId: str
    side: str
    qty: str
    orderPrice: str
    orderType: str
    execType: str
    closedSize: str
    cumEntryValue: str
    avgEntryPrice: str
    cumExitValue: str
    avgExitPrice: str
    closedPnl: str
    fillCount: str
    leverage: str
    createdTime: str
    updatedTime: str


class Result(BaseModel):
    list: List[TradeOrder]
    category: str
    nextPageCursor: str


class JsonResponse(BaseModel):
    retCode: int
    retMsg: str
    result: Result
    retExtInfo: Dict
    time: int


def get_closed_pnl(api_key, secret, days=1) -> JsonResponse:
    print(f"{api_key=} {secret=}")
    session = HTTP(
        testnet=False,
        api_key=api_key,
        api_secret=secret,
    )

    from datetime import datetime, timedelta

    now = datetime.now()
    yesterday = now - timedelta(days=days)

    start_time = yesterday.replace(hour=0, minute=0, second=0, microsecond=0)
    end_time = now.replace(hour=23, minute=59, second=59, microsecond=999999)

    start_timestamp = int(start_time.timestamp() * 1000)  # Convert to milliseconds
    end_timestamp = int(end_time.timestamp() * 1000)  # Convert to milliseconds

    r = session.get_closed_pnl(
        category="inverse",
        symbol="ETHUSD",
        startTime=start_timestamp,
        endTime=end_timestamp,
    )

    return JsonResponse.parse_obj(r)


def main() -> int:
    try:
        with SessionLocal() as db:
            active_api_keys = (
                db.query(ApiKey)
                .filter(ApiKey.status == ApiKeyStatusEnum.ACTIVE.value)
                .all()
            )
            if not active_api_keys:
                print(f"There is no active api key")
                return 1
            for ak in active_api_keys:
                if ak.status == ApiKeyStatusEnum.ACTIVE.value:
                    resp = get_closed_pnl(api_key=ak.api_key, secret=ak.secret, days=1)
                    all = [
                        ClosedPnl(user_id=ak.user_id, api_key_id=ak.id, **r.dict())
                        for r in resp.result.list
                    ]
                    print(
                        f"There is {len(all)} new register for the user {ak.user.username}"
                    )
                    db.add_all(all)
                    db.commit()
    except Exception as e:
        print(f"Something Wrong in main function. {e}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
